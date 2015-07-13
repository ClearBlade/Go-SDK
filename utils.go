package GoSDK

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	mqtt "github.com/clearblade/mqtt_parsing"
	"github.com/clearblade/mqttclient"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	CB_ADDR            = "https://rtp.clearblade.com"
	CB_MSG_ADDR        = "rtp.clearblade.com:1883"
	_HEADER_KEY_KEY    = "ClearBlade-SystemKey"
	_HEADER_SECRET_KEY = "ClearBlade-SystemSecret"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
}

//Client is a convience interface for API consumers, if they want to use the same functions for both
//Dev Users and unprivleged users, such as tiny helper functions. please use this with care
type Client interface {
	//bookkeeping calls
	Authenticate() error
	Logout() error

	//data calls
	InsertData(string, interface{}) error
	UpdateData(string, *Query, map[string]interface{}) error

	GetData(string, *Query) (map[string]interface{}, error)
	GetDataByName(string, *Query) (map[string]interface{}, error)
	GetDataByKeyAndName(string, string, *Query) (map[string]interface{}, error)

	DeleteData(string, *Query) error

	//mqtt calls
	InitializeMQTT(string, string, int) error
	ConnectMQTT(*tls.Config, *LastWillPacket) error
	Publish(string, []byte, int) error
	Subscribe(string, int) (<-chan *mqtt.Publish, error)
	Unsubscribe(string) error
	Disconnect() error
}

//cbClient will supply various information that differs between privleged and unprivleged users
//this interface is meant to be unexported
type cbClient interface {
	credentials() ([][]string, error) //the inner slice is a tuple of "Header":"Value"
	preamble() string
	setToken(string)
	getToken() string
	getSystemInfo() (string, string)
	getMessageId() uint16
}

type UserClient struct {
	UserToken    string
	mrand        *rand.Rand
	MQTTClient   *mqttclient.Client
	SystemKey    string
	SystemSecret string
	Email        string
	Password     string
}

type DevClient struct {
	DevToken   string
	mrand      *rand.Rand
	MQTTClient *mqttclient.Client
	Email      string
	Password   string
}

type CbReq struct {
	Body        interface{}
	Method      string
	Endpoint    string
	QueryString string
	Headers     map[string][]string
}

type CbResp struct {
	Body       interface{}
	StatusCode int
}

func NewUserClient(systemkey, systemsecret, email, password string) *UserClient {
	return &UserClient{
		UserToken:    "",
		mrand:        rand.New(rand.NewSource(time.Now().UnixNano())),
		MQTTClient:   nil,
		SystemSecret: systemsecret,
		SystemKey:    systemkey,
		Email:        email,
		Password:     password,
	}
}

func NewDevClient(email, password string) *DevClient {
	return &DevClient{
		DevToken:   "",
		mrand:      rand.New(rand.NewSource(time.Now().UnixNano())),
		MQTTClient: nil,
		Email:      email,
		Password:   password,
	}
}

func (u *UserClient) Authenticate() error {
	return authenticate(u, u.Email, u.Password)
}

func (d *DevClient) Authenticate() error {
	return authenticate(d, d.Email, d.Password)
}

func (u *UserClient) Register(username, password string) error {
	_, err := register(u, username, password, "", "", "")
	return err
}

func (u *UserClient) RegisterUser(username, password string) (map[string]interface{}, error) {
	resp, err := register(u, username, password, "", "", "")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DevClient) Register(username, password, fname, lname, org string) error {
	_, err := register(d, username, password, fname, lname, org)
	return err
}

func (d *DevClient) RegisterUser(username, password, fname, lname, org string) (map[string]interface{}, error) {
	resp, err := register(d, username, password, fname, lname, org)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *UserClient) Logout() error {
	return logout(u)
}

func (d *DevClient) Logout() error {
	return logout(d)
}

//Below are some shared functions

func authenticate(c cbClient, username, password string) error {
	var creds [][]string
	switch c.(type) {
	case *UserClient:
		var err error
		creds, err = c.credentials()
		if err != nil {
			return err
		}
	}

	fmt.Printf("AUTHING WITH %s:%s\n", username, password)
	resp, err := post(c.preamble()+"/auth", map[string]interface{}{
		"email":    username,
		"password": password,
	}, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error in authenticating, Status Code: %d, %v\n", resp.StatusCode, resp.Body)
	}

	var token string = ""
	switch c.(type) {
	case *UserClient:
		token = resp.Body.(map[string]interface{})["user_token"].(string)
	case *DevClient:
		token = resp.Body.(map[string]interface{})["dev_token"].(string)
	}
	if token == "" {
		return fmt.Errorf("Token not present i response from platform %+v", resp.Body)
	}
	c.setToken(token)
	return nil
}

func register(c cbClient, username, password, fname, lname, org string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"email":    username,
		"password": password,
	}

	creds := [][]string{}
	switch c.(type) {
	case *DevClient:
		payload["fname"] = fname
		payload["lname"] = lname
		payload["org"] = org
	case *UserClient:
		var err error
		creds, err = c.credentials()
		if err != nil {
			return nil, err
		}
	}

	resp, err := post(c.preamble()+"/reg", payload, creds, nil)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Status code: %d, Error in authenticating, %v\n", resp.StatusCode, resp.Body)
	}
	var token string = ""
	switch c.(type) {
	case *UserClient:
		token = resp.Body.(map[string]interface{})["user_token"].(string)
	case *DevClient:
		token = resp.Body.(map[string]interface{})["dev_token"].(string)
	}
	if token == "" {
		return nil, fmt.Errorf("Token not present i response from platform %+v", resp.Body)
	}
	c.setToken(token)
	return resp.Body.(map[string]interface{}), nil
}

func logout(c cbClient) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	resp, err := post(c.preamble()+"/logout", nil, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error in authenticating %v\n", resp.Body)
	}
	return nil
}

func do(r *CbReq, creds [][]string) (*CbResp, error) {
	var bodyToSend *bytes.Buffer
	if r.Body != nil {
		b, jsonErr := json.Marshal(r.Body)
		if jsonErr != nil {
			return nil, fmt.Errorf("JSON Encoding Error: %v", jsonErr)
		}
		bodyToSend = bytes.NewBuffer(b)
	} else {
		bodyToSend = nil
	}
	url := CB_ADDR + r.Endpoint
	if r.QueryString != "" {
		url += "?" + r.QueryString
	}
	var req *http.Request
	var reqErr error
	if bodyToSend != nil {
		req, reqErr = http.NewRequest(r.Method, url, bodyToSend)
	} else {
		req, reqErr = http.NewRequest(r.Method, url, nil)
	}
	if reqErr != nil {
		return nil, fmt.Errorf("Request Creation Error: %v", reqErr)
	}
	if r.Headers != nil {
		for hed, val := range r.Headers {
			for _, vv := range val {
				req.Header.Add(hed, vv)
			}
		}
	}
	for _, c := range creds {
		if len(c) != 2 {
			return nil, fmt.Errorf("Request Creation Error: Invalid credential header supplied")
		}
		req.Header.Add(c[0], c[1])
	}

	cli := &http.Client{Transport: tr}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error Making Request: %v", err)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("Error Reading Response Body: %v", readErr)
	}
	var d interface{}
	if len(body) == 0 {
		return &CbResp{
			Body:       nil,
			StatusCode: resp.StatusCode,
		}, nil
	}
	buf := bytes.NewBuffer(body)
	dec := json.NewDecoder(buf)
	decErr := dec.Decode(&d)
	var bod interface{}
	if decErr != nil {
		//		return nil, fmt.Errorf("JSON Decoding Error: %v\n With Body: %v\n", decErr, string(body))
		bod = string(body)
	}
	switch d.(type) {
	case []interface{}:
		bod = d
	case map[string]interface{}:
		bod = d
	default:
		bod = string(body)
	}
	return &CbResp{
		Body:       bod,
		StatusCode: resp.StatusCode,
	}, nil
}

//standard http verbs

func get(endpoint string, query map[string]string, creds [][]string, headers map[string][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        nil,
		Method:      "GET",
		Endpoint:    endpoint,
		QueryString: query_to_string(query),
		Headers:     headers,
	}
	return do(req, creds)
}

func post(endpoint string, body interface{}, creds [][]string, headers map[string][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        body,
		Method:      "POST",
		Endpoint:    endpoint,
		QueryString: "",
		Headers:     headers,
	}
	return do(req, creds)
}

func put(endpoint string, body interface{}, heads [][]string, headers map[string][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        body,
		Method:      "PUT",
		Endpoint:    endpoint,
		QueryString: "",
		Headers:     headers,
	}
	return do(req, heads)
}

func delete(endpoint string, query map[string]string, heds [][]string, headers map[string][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        nil,
		Method:      "DELETE",
		Endpoint:    endpoint,
		Headers:     headers,
		QueryString: query_to_string(query),
	}
	return do(req, heds)
}

func query_to_string(query map[string]string) string {
	qryStr := ""
	for k, v := range query {
		qryStr += k + "=" + v + "&"
	}
	return strings.TrimSuffix(qryStr, "&")
}
