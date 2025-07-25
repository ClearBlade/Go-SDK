package GoSDK

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	cbErr "github.com/clearblade/go-utils/errors"
	mqttTypes "github.com/clearblade/mqtt_parsing"
	mqtt "github.com/clearblade/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
)

var (
	//CB_ADDR is the address of the ClearBlade Platform you are speaking with
	CB_ADDR = "https://platform.clearblade.com"
	//CB_MSG_ADDR is the messaging address you wish to speak to
	CB_MSG_ADDR      = "platform.clearblade.com:1883"
	CB_MSG_AUTH_ADDR = "platform.clearblade.com:8905"

	MTLS_PORT = "444"

	_HEADER_KEY_KEY    = "ClearBlade-SystemKey"
	_HEADER_SECRET_KEY = "ClearBlade-SystemSecret"
)

var tr = &http.Transport{
	// TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

const (
	createDevUser = iota
	createUser
)

// Client is a convience interface for API consumers, if they want to use the same functions for both developer users and unprivleged users
type Client interface {
	//session bookkeeping calls
	Authenticate() (*AuthResponse, error)
	Logout() error

	//data calls
	CreateData(string, interface{}) ([]interface{}, error)
	CreateDataByName(string, string, interface{}) ([]interface{}, error)
	InsertData(string, interface{}) error
	UpdateData(string, *Query, map[string]interface{}) error
	UpdateDataByName(string, string, *Query, map[string]interface{}) (UpdateResponse, error)
	GetData(string, *Query) (map[string]interface{}, error)
	GetDataByName(string, *Query) (map[string]interface{}, error)
	GetDataByNameWithSystemKey(string, string, *Query) (map[string]interface{}, error)
	DeleteData(string, *Query) error
	GetItemCount(string) (int, error)
	GetDataTotal(string, *Query) (map[string]interface{}, error)
	GetColumns(string, string, string) ([]interface{}, error)

	//mqtt calls
	SetMqttClient(MqttClient)
	InitializeMQTT(string, string, int, *tls.Config, *LastWillPacket) error
	Publish(string, []byte, int) error
	Subscribe(string, int) (<-chan *mqttTypes.Publish, error)
	Unsubscribe(string) error
	Disconnect() error

	// Device calls
	GetDevices(string, *Query) ([]interface{}, error)
	GetDevice(string, string) (map[string]interface{}, error)
	CreateDevice(string, string, map[string]interface{}) (map[string]interface{}, error)
	UpdateDevice(string, string, map[string]interface{}) (map[string]interface{}, error)
	DeleteDevice(string, string) error
	UpdateDevices(string, *Query, map[string]interface{}) ([]interface{}, error)
	DeleteDevices(string, *Query) error

	// Adaptor calls
	GetAdaptors(string) ([]interface{}, error)
	GetAdaptor(string, string) (map[string]interface{}, error)
	CreateAdaptor(string, string, map[string]interface{}) (map[string]interface{}, error)
	UpdateAdaptor(string, string, map[string]interface{}) (map[string]interface{}, error)
	DeleteAdaptor(string, string) error
	DeployAdaptor(string, string, map[string]interface{}) (map[string]interface{}, error)
	ControlAdaptor(string, string, map[string]interface{}) (map[string]interface{}, error)

	// Adaptor File calls
	GetAdaptorFiles(string, string) ([]interface{}, error)
	GetAdaptorFile(string, string, string) (map[string]interface{}, error)
	CreateAdaptorFile(string, string, string, map[string]interface{}) (map[string]interface{}, error)
	UpdateAdaptorFile(string, string, string, map[string]interface{}) (map[string]interface{}, error)
	DeleteAdaptorFile(string, string, string) error
}

// interface assertions
var (
	_ Client = (*UserClient)(nil)
	_ Client = (*DeviceClient)(nil)
	_ Client = (*DevClient)(nil)
)

type MqttClient interface {
	mqtt.Client
}

// cbClient will supply various information that differs between privleged and unprivleged users
// this interface is meant to be unexported
type cbClient interface {
	credentials() ([][]string, error) //the inner slice is a tuple of "Header":"Value"
	preamble() string
	setToken(string)
	getToken() string
	setRefreshToken(map[string]interface{})
	getRefreshToken() string
	setExpiresAt(float64)
	getExpiresAt() float64
	getSystemInfo() (string, string)
	getMessageId() uint16
	getHttpAddr() string
	getMqttAddr() string
	getMTLSPort() string
	getEdgeProxy() *EdgeProxy
}

// receiver for methods that can be shared between users/devs/devices
type client struct{}

// UserClient is the type for users
type UserClient struct {
	client
	UserToken    string
	RefreshToken string
	ExpiresAt    float64
	MQTTClient   MqttClient
	SystemKey    string
	SystemSecret string
	Email        string
	Password     string
	HttpAddr     string
	MqttAddr     string
	MqttAuthAddr string
	MTLSPort     string
	edgeProxy    *EdgeProxy
}

type DeviceClient struct {
	client
	DeviceName   string
	ActiveKey    string
	DeviceToken  string
	RefreshToken string
	ExpiresAt    float64
	MQTTClient   MqttClient
	SystemKey    string
	SystemSecret string
	HttpAddr     string
	MqttAddr     string
	MqttAuthAddr string
	MTLSPort     string
	edgeProxy    *EdgeProxy
	IsMTLS       bool
	Cert         string
	Key          string
	DetailsInCN  bool
}

// DevClient is the type for developers
type DevClient struct {
	client
	DevToken     string
	RefreshToken string
	ExpiresAt    float64
	MQTTClient   MqttClient
	Email        string
	Password     string
	HttpAddr     string
	MqttAddr     string
	BrokerWsAddr string
	MqttAuthAddr string
	MTLSPort     string
	edgeProxy    *EdgeProxy
}

type EdgeProxy struct {
	SystemKey string
	EdgeName  string
}

// CbReq is a wrapper around an HTTP request
type CbReq struct {
	Body        interface{}
	Method      string
	Endpoint    string
	QueryString string
	Headers     map[string][]string
	HttpAddr    string
	MqttAddr    string
	Transport   *http.Transport
	IsMTLS      bool
}

func (r *CbReq) setupForMTLS(client *DeviceClient) error {
	c, err := tls.X509KeyPair([]byte(client.Cert), []byte(client.Key))
	if err != nil {
		return fmt.Errorf("Error loading X509 Key Pair: %v", err)
	}
	r.IsMTLS = true
	r.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{c},
			InsecureSkipVerify: true,
		},
	}
	return nil
}

// CbResp is a wrapper around an HTTP response
type CbResp struct {
	Body       interface{}
	StatusCode int
}

type AuthResponse struct {
	DevResponse *DevAuthResponse
}

type DevAuthResponse struct {
	DevToken          string `json:"dev_token"`
	IsTwoFactor       bool   `json:"is_two_factor"`
	NextStepURL       string `json:"next_step_url"`
	IntermediateToken string `json:"intermediate_token"`
	TwoFactorMethod   string `json:"two_factor_method"`
	OtpID             string `json:"otp_id"`
	OtpIssued         string `json:"otp_issued"`
}

func (u *UserClient) getHttpAddr() string {
	return u.HttpAddr
}

func (u *UserClient) getMTLSPort() string {
	return u.MTLSPort
}

func (d *DevClient) getHttpAddr() string {
	return d.HttpAddr
}

func (d *DevClient) getMTLSPort() string {
	return d.MTLSPort
}

func (u *UserClient) getMqttAddr() string {
	return u.MqttAddr
}

func (d *DevClient) getMqttAddr() string {
	return d.MqttAddr
}

func (u *UserClient) getEdgeProxy() *EdgeProxy {
	return u.edgeProxy
}

func (d *DevClient) getEdgeProxy() *EdgeProxy {
	return d.edgeProxy
}

func (d *DeviceClient) getEdgeProxy() *EdgeProxy {
	return d.edgeProxy
}

func (u *UserClient) SetMqttClient(c MqttClient) {
	u.MQTTClient = c
}

func (d *DevClient) SetMqttClient(c MqttClient) {
	d.MQTTClient = c
}

func (d *DeviceClient) SetMqttClient(c MqttClient) {
	d.MQTTClient = c
}

func NewDeviceClient(systemkey, systemsecret, deviceName, activeKey string) *DeviceClient {
	return &DeviceClient{
		DeviceName:   deviceName,
		DeviceToken:  "",
		RefreshToken: "",
		ActiveKey:    activeKey,
		MQTTClient:   nil,
		SystemKey:    systemkey,
		SystemSecret: systemsecret,
		HttpAddr:     CB_ADDR,
		MqttAddr:     CB_MSG_ADDR,
		MqttAuthAddr: CB_MSG_AUTH_ADDR,
		MTLSPort:     MTLS_PORT,
	}
}

func NewDeviceMTLSClient(systemkey, deviceName, cert, key string, detailsInCommonName bool) *DeviceClient {
	return &DeviceClient{
		DeviceName:   deviceName,
		DeviceToken:  "",
		RefreshToken: "",
		MQTTClient:   nil,
		SystemKey:    systemkey,
		HttpAddr:     CB_ADDR,
		MqttAddr:     CB_MSG_ADDR,
		MqttAuthAddr: CB_MSG_AUTH_ADDR,
		MTLSPort:     MTLS_PORT,
		IsMTLS:       true,
		Cert:         cert,
		Key:          key,
		DetailsInCN:  detailsInCommonName,
	}
}

// NewUserClient allocates a new UserClient struct
func NewUserClient(systemkey, systemsecret, email, password string) *UserClient {
	return &UserClient{
		UserToken:    "",
		RefreshToken: "",
		MQTTClient:   nil,
		SystemSecret: systemsecret,
		SystemKey:    systemkey,
		Email:        email,
		Password:     password,
		HttpAddr:     CB_ADDR,
		MqttAddr:     CB_MSG_ADDR,
		MqttAuthAddr: CB_MSG_AUTH_ADDR,
		MTLSPort:     MTLS_PORT,
	}
}

// NewDevClient allocates a new DevClient struct
func NewDevClient(email, password string) *DevClient {
	return &DevClient{
		DevToken:     "",
		RefreshToken: "",
		MQTTClient:   nil,
		Email:        email,
		Password:     password,
		HttpAddr:     CB_ADDR,
		MqttAddr:     CB_MSG_ADDR,
		MqttAuthAddr: CB_MSG_AUTH_ADDR,
		MTLSPort:     MTLS_PORT,
	}
}

func NewDevClientWithToken(token, email string) *DevClient {
	return &DevClient{
		DevToken:     token,
		RefreshToken: "",
		MQTTClient:   nil,
		Email:        email,
		Password:     "",
		HttpAddr:     CB_ADDR,
		MqttAddr:     CB_MSG_ADDR,
		MqttAuthAddr: CB_MSG_AUTH_ADDR,
		MTLSPort:     MTLS_PORT,
	}
}

func NewRefreshUserClientWithAddrs(httpAddr, mqttAddr, systemKey, systemSecret, refreshToken, accessToken string) *UserClient {
	return &UserClient{
		UserToken:    accessToken,
		RefreshToken: refreshToken,
		MQTTClient:   nil,
		SystemSecret: systemSecret,
		SystemKey:    systemKey,
		HttpAddr:     httpAddr,
		MqttAddr:     mqttAddr,
		MqttAuthAddr: CB_MSG_AUTH_ADDR,
		MTLSPort:     MTLS_PORT,
	}
}

func NewUserClientWithAddrs(httpAddr, mqttAddr, systemKey, systemSecret, email, password string) *UserClient {
	return &UserClient{
		UserToken:    "",
		RefreshToken: "",
		MQTTClient:   nil,
		SystemSecret: systemSecret,
		SystemKey:    systemKey,
		Email:        email,
		Password:     password,
		HttpAddr:     httpAddr,
		MqttAddr:     mqttAddr,
		MqttAuthAddr: CB_MSG_AUTH_ADDR,
		MTLSPort:     MTLS_PORT,
	}
}

func NewUserClientWithAddrs2(httpAddr, mqttAddr, mqttAuthAddr, systemKey, systemSecret, email, password string) *UserClient {
	return &UserClient{
		UserToken:    "",
		RefreshToken: "",
		MQTTClient:   nil,
		SystemSecret: systemSecret,
		SystemKey:    systemKey,
		Email:        email,
		Password:     password,
		HttpAddr:     httpAddr,
		MqttAddr:     mqttAddr,
		MqttAuthAddr: mqttAuthAddr,
		MTLSPort:     MTLS_PORT,
	}
}

func NewDevClientWithAddrs(httpAddr, mqttAddr, email, password string) *DevClient {
	return &DevClient{
		DevToken:     "",
		RefreshToken: "",
		MQTTClient:   nil,
		Email:        email,
		Password:     password,
		HttpAddr:     httpAddr,
		MqttAddr:     mqttAddr,
		MqttAuthAddr: CB_MSG_AUTH_ADDR,
		MTLSPort:     MTLS_PORT,
	}
}

func NewDevClientWithWebsocket(httpAddr, mqttAddr, wsAddr, email, password string) *DevClient {
	c := NewDevClientWithAddrs(httpAddr, mqttAddr, email, password)
	c.BrokerWsAddr = wsAddr
	return c
}

func NewDevClientWithTokenAndAddrs(httpAddr, mqttAddr, token, email string) *DevClient {
	return &DevClient{
		DevToken:     token,
		RefreshToken: "",
		MQTTClient:   nil,
		Email:        email,
		Password:     "",
		HttpAddr:     httpAddr,
		MqttAddr:     mqttAddr,
		MqttAuthAddr: CB_MSG_AUTH_ADDR,
		MTLSPort:     MTLS_PORT,
	}
}

func NewDeviceClientWithAddrs(httpAddr, mqttAddr, systemkey, systemsecret, deviceName, activeKey string) *DeviceClient {
	return &DeviceClient{
		DeviceName:   deviceName,
		DeviceToken:  "",
		RefreshToken: "",
		ActiveKey:    activeKey,
		MQTTClient:   nil,
		SystemKey:    systemkey,
		SystemSecret: systemsecret,
		HttpAddr:     httpAddr,
		MqttAddr:     mqttAddr,
		MqttAuthAddr: CB_MSG_AUTH_ADDR,
		MTLSPort:     MTLS_PORT,
	}
}

func NewDeviceClientWithServiceAccountAndAddrs(httpAddr, mqttAddr, systemkey, systemsecret, deviceName, token string) *DeviceClient {
	return &DeviceClient{
		DeviceName:   deviceName,
		DeviceToken:  token,
		ActiveKey:    "",
		MQTTClient:   nil,
		SystemKey:    systemkey,
		SystemSecret: systemsecret,
		HttpAddr:     httpAddr,
		MqttAddr:     mqttAddr,
		MTLSPort:     MTLS_PORT,
	}
}

func NewUserClientWithServiceAccountAndAddrs(httpAddr, mqttAddr, systemkey, systemsecret, email, token string) *UserClient {
	return &UserClient{
		UserToken:    token,
		MQTTClient:   nil,
		SystemSecret: systemsecret,
		SystemKey:    systemkey,
		Email:        email,
		Password:     "",
		HttpAddr:     httpAddr,
		MqttAddr:     mqttAddr,
		MTLSPort:     MTLS_PORT,
	}
}

func NewEdgeProxyDevClient(email, password, systemKey, edgeName string) (*DevClient, error) {
	d := NewDevClient(email, password)
	if err := d.startProxyToEdge(systemKey, edgeName); err != nil {
		return nil, err
	}
	return d, nil
}
func NewEdgeProxyUserClient(email, password, systemKey, systemSecret, edgeName string) (*UserClient, error) {
	u := NewUserClient(systemKey, systemSecret, email, password)
	if err := u.startProxyToEdge(systemKey, edgeName); err != nil {
		return nil, err
	}
	return u, nil
}
func NewEdgeProxyDeviceClient(systemkey, systemsecret, deviceName, activeKey, edgeName string) (*DeviceClient, error) {
	d := NewDeviceClient(systemkey, systemsecret, deviceName, activeKey)
	if err := d.startProxyToEdge(systemkey, edgeName); err != nil {
		return nil, err
	}
	return d, nil
}

func (u *UserClient) startProxyToEdge(systemKey, edgeName string) error {
	if systemKey == "" || edgeName == "" {
		return fmt.Errorf("systemKey and edgeName required")
	}
	u.edgeProxy = &EdgeProxy{systemKey, edgeName}
	return nil
}
func (u *UserClient) stopProxyToEdge() error {
	if u.edgeProxy == nil {
		return fmt.Errorf("Requests are not being proxied to edge")
	}
	u.edgeProxy = nil
	return nil
}
func (d *DevClient) startProxyToEdge(systemKey, edgeName string) error {
	if systemKey == "" || edgeName == "" {
		return fmt.Errorf("systemKey and edgeName required")
	}
	d.edgeProxy = &EdgeProxy{systemKey, edgeName}
	return nil
}
func (d *DevClient) stopProxyToEdge() error {
	if d.edgeProxy == nil {
		return fmt.Errorf("No edge proxy active")
	}
	d.edgeProxy = nil
	return nil
}
func (d *DeviceClient) startProxyToEdge(systemKey, edgeName string) error {
	if systemKey == "" || edgeName == "" {
		return fmt.Errorf("systemKey and edgeName required")
	}
	d.edgeProxy = &EdgeProxy{systemKey, edgeName}
	return nil
}
func (d *DeviceClient) stopProxyToEdge() error {
	if d.edgeProxy == nil {
		return fmt.Errorf("No edge proxy active")
	}
	d.edgeProxy = nil
	return nil
}

// Authenticate retrieves a token from the specified Clearblade Platform
func (u *UserClient) Authenticate() (*AuthResponse, error) {
	if err := authenticate(u, u.Email, u.Password, map[string]string{}); err != nil {
		return nil, err
	}
	return nil, nil
}

func (u *UserClient) AuthenticateWithOptions(opts map[string]string) (*AuthResponse, error) {
	if err := authenticate(u, u.Email, u.Password, opts); err != nil {
		return nil, err
	}
	return nil, nil
}

func (u *UserClient) RefreshAuthentication() error {
	if err := refreshAuthentication(u); err != nil {
		return err
	}
	return nil
}

func (u *UserClient) AuthAnon() error {
	return authAnon(u)
}

// Authenticate retrieves a token from the specified Clearblade Platform
func (d *DevClient) Authenticate() (*AuthResponse, error) {
	var creds [][]string
	resp, err := post(d, d.preamble()+"/auth", map[string]interface{}{
		"email":    d.Email,
		"password": d.Password,
	}, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, cbErr.CreateResponseFromMap(resp.Body)
	}

	body := resp.Body.(map[string]interface{})
	var devAuthResp *DevAuthResponse

	if val, ok := body["is_two_factor"]; ok && val.(bool) {
		devAuthResp = &DevAuthResponse{
			DevToken:          body["dev_token"].(string),
			IsTwoFactor:       body["is_two_factor"].(bool),
			NextStepURL:       body["next_step_url"].(string),
			IntermediateToken: body["intermediate_token"].(string),
			TwoFactorMethod:   body["two_factor_method"].(string),
			OtpID:             body["otp_id"].(string),
			OtpIssued:         body["otp_issued"].(string),
		}
	} else {
		devAuthResp = &DevAuthResponse{
			DevToken:          body["dev_token"].(string),
			IsTwoFactor:       false,
			NextStepURL:       "",
			IntermediateToken: "",
			TwoFactorMethod:   "",
			OtpID:             "",
			OtpIssued:         "",
		}
	}

	var token string
	if devAuthResp.IsTwoFactor {
		token = devAuthResp.IntermediateToken
	} else {
		token = devAuthResp.DevToken
	}

	if token == "" {
		return nil, fmt.Errorf("Token not present in response from platform %+v", resp.Body)
	}
	d.setToken(token)
	return &AuthResponse{
		DevResponse: devAuthResp,
	}, nil
}

type VerifyAuthenticationParams struct {
	Code            string `json:"code"`
	TwoFactorMethod string `json:"two_factor_method"`
	OtpID           string `json:"otp_id"`
	OtpIssued       string `json:"otp_issued"`
}

func (d *DevClient) VerifyAuthentication(verifyParams VerifyAuthenticationParams) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := post(d, d.preamble()+"/auth/verify", map[string]interface{}{
		"code":              verifyParams.Code,
		"two_factor_method": verifyParams.TwoFactorMethod,
		"otp_id":            verifyParams.OtpID,
		"otp_issued":        verifyParams.OtpIssued,
	}, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return cbErr.CreateResponseFromMap(resp.Body)
	}

	body := resp.Body.(map[string]interface{})
	d.setToken(body["dev_token"].(string))
	return nil
}

// Register creates a new user
func (u *UserClient) Register(username, password string) error {
	if u.UserToken == "" {
		return fmt.Errorf("Must be logged in to create users")
	}
	_, err := register(u, createUser, username, password, u.SystemKey, u.SystemSecret, "", "", "", "")
	return err
}

// RegisterUser creates a new user, returning the body of the response.
func (u *UserClient) RegisterUser(username, password string) (map[string]interface{}, error) {
	if u.UserToken == "" {
		return nil, fmt.Errorf("Must be logged in to create users")
	}
	resp, err := register(u, createUser, username, password, u.SystemKey, u.SystemSecret, "", "", "", "")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Registers a new developer
func (d *DevClient) Register(username, password, fname, lname, org string) error {
	resp, err := register(d, createDevUser, username, password, "", "", fname, lname, org, "")
	if err != nil {
		return err
	} else {
		d.DevToken = resp["dev_token"].(string)
		return nil
	}
}

func (d *DevClient) RegisterNewUser(username, password, systemkey, systemsecret string) (map[string]interface{}, error) {
	if d.DevToken == "" {
		return nil, fmt.Errorf("Must authenticate first")
	}
	return register(d, createUser, username, password, systemkey, systemsecret, "", "", "", "")

}

// Register creates a new developer user
func (d *DevClient) RegisterDevUser(username, password, fname, lname, org string) (map[string]interface{}, error) {
	resp, err := register(d, createDevUser, username, password, "", "", fname, lname, org, "")
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Register creates a new developer user
func (d *DevClient) RegisterDevUserWithKey(username, password, fname, lname, org, key string) (map[string]interface{}, error) {
	resp, err := register(d, createDevUser, username, password, "", "", fname, lname, org, key)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Logout ends the session
func (u *UserClient) Logout() error {
	return logout(u)
}

// Logout ends the session
func (d *DevClient) Logout() error {
	return logout(d)
}

// Check Auth of User
func (d *UserClient) CheckAuth() error {
	return checkAuth(d)
}

// Check Auth of Developer
func (d *DevClient) CheckAuth() error {
	return checkAuth(d)
}

func checkAuth(c cbClient) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	//log.Println("Checking user auth")
	resp, err := post(c, c.preamble()+"/checkauth", nil, creds, nil)
	if err != nil {
		return err
	}
	body := resp.Body.(map[string]interface{})
	if body["is_authenticated"] != nil && body["is_authenticated"].(bool) {
		return nil
	}
	if resp.StatusCode != 200 {
		return cbErr.CreateResponseFromMap(resp.Body)
	}
	return nil
}

//Below are some shared functions

func authenticate(c cbClient, username, password string, opts map[string]string) error {
	var creds [][]string
	switch c.(type) {
	case *UserClient:
		var err error
		creds, err = c.credentials()
		if err != nil {
			return err
		}
	case *DevClient:
	}
	if opts == nil {
		opts = map[string]string{}
	}
	opts["email"] = username
	opts["password"] = password
	resp, err := post(c, c.preamble()+"/auth", opts, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return cbErr.CreateResponseFromMap(resp.Body)
	}

	var token string = ""
	respBody := resp.Body.(map[string]interface{})
	switch c.(type) {
	case *UserClient:
		token = respBody["user_token"].(string)
	case *DevClient:
		token = respBody["dev_token"].(string)
	}
	if token == "" {
		return fmt.Errorf("Token not present i response from platform %+v", resp.Body)
	}
	c.setToken(token)
	if _, ok := respBody["refresh_token"]; ok {
		c.setRefreshToken(respBody)
	}
	if _, ok := respBody["expires_at"]; ok {
		c.setExpiresAt(respBody["expires_at"].(float64))
	}
	return nil
}

func refreshAuthentication(c *UserClient) error {
	var creds [][]string
	var err error
	creds, err = c.credentials()
	if err != nil {
		return err
	}

	resp, err := post(c, c.preamble()+"/auth", map[string]interface{}{
		"refresh_token": c.getRefreshToken(),
		"access_token":  c.getToken(),
		"grant_type":    "refresh_token",
	}, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return cbErr.CreateResponseFromMap(resp.Body)
	}

	respBody := resp.Body.(map[string]interface{})
	c.setToken(respBody["user_token"].(string))
	c.setRefreshToken(respBody)
	c.setExpiresAt(respBody["expires_at"].(float64))
	return nil
}

func authAnon(c cbClient) error {
	creds, err := c.credentials()
	if err != nil {
		return fmt.Errorf("Invalid client: %+s", err.Error())
	}
	resp, err := post(c, c.preamble()+"/anon", nil, creds, nil)
	if err != nil {
		return fmt.Errorf("Error retrieving anon user token: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		return cbErr.CreateResponseFromMap(resp.Body)
	}
	token := resp.Body.(map[string]interface{})["user_token"].(string)
	if token == "" {
		return fmt.Errorf("Token not present in response from platform %+v", resp.Body)
	}
	c.setToken(token)
	return nil
}

func register(c cbClient, kind int, username, password, syskey, syssec, fname, lname, org, key string) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"email":    username,
		"password": password,
	}
	var endpoint string
	headers := make(map[string][]string)
	var creds [][]string
	switch kind {
	case createDevUser:
		endpoint = "/admin/reg"
		payload["fname"] = fname
		payload["lname"] = lname
		payload["org"] = org
		if key != "" {
			payload["key"] = key
		}
	case createUser:
		switch c.(type) {
		case *DevClient:
			if syskey == "" {
				return nil, fmt.Errorf("System key required")
			}
			endpoint = fmt.Sprintf("/admin/user/%s", syskey)
		case *UserClient:
			if syskey == "" {
				return nil, fmt.Errorf("System key required")
			}
			if syssec == "" {
				return nil, fmt.Errorf("System secret required")
			}
			endpoint = "/api/v/1/user/reg"
			headers["Clearblade-Systemkey"] = []string{syskey}
			headers["Clearblade-Systemsecret"] = []string{syssec}
		default:
			return nil, fmt.Errorf("unreachable code detected")
		}
		var err error
		creds, err = c.credentials()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Cannot create that kind of user")
	}
	resp, err := post(c, endpoint, payload, creds, headers)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, cbErr.CreateResponseFromMap(resp.Body)
	}
	var token string = ""
	switch kind {
	case createDevUser:
		token = resp.Body.(map[string]interface{})["dev_token"].(string)
	case createUser:
		token = resp.Body.(map[string]interface{})["user_id"].(string)
	}
	c.setExpiresAt(resp.Body.(map[string]interface{})["expires_at"].(float64))
	c.setRefreshToken(resp.Body.(map[string]interface{}))

	if token == "" {
		return nil, fmt.Errorf("Token not present in response from platform %+v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func logout(c cbClient) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	resp, err := post(c, c.preamble()+"/logout", nil, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return cbErr.CreateResponseFromMap(resp.Body)
	}
	return nil
}

func do(c cbClient, r *CbReq, creds [][]string) (*CbResp, error) {
	checkForEdgeProxy(c, r)
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
	url := c.getHttpAddr() + r.Endpoint
	if r.IsMTLS {
		url = c.getHttpAddr() + ":" + c.getMTLSPort() + r.Endpoint
	}
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
		return nil, fmt.Errorf("Request Creation Error: %s", reqErr)
	}
	req.Close = true
	for hed, val := range r.Headers {
		for _, vv := range val {
			req.Header.Add(hed, vv)
		}
	}
	for _, c := range creds {
		if len(c) != 2 {
			return nil, fmt.Errorf("Request Creation Error: Invalid credential header supplied")
		}
		req.Header.Add(c[0], c[1])
	}
	cli := &http.Client{
		Transport: tr,
		Timeout:   time.Minute * 5,
	}
	if r.Transport != nil {
		cli.Transport = r.Transport
	}
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

func get(c cbClient, endpoint string, query map[string]string, creds [][]string, headers map[string][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        nil,
		Method:      "GET",
		Endpoint:    endpoint,
		QueryString: query_to_string(query),
		Headers:     headers,
	}
	d, ok := c.(*DeviceClient)
	if ok && d.IsMTLS {
		if err := req.setupForMTLS(d); err != nil {
			return nil, fmt.Errorf("Error setting up MTLS: %v", err)
		}
	}
	return do(c, req, creds)
}

func post(c cbClient, endpoint string, body interface{}, creds [][]string, headers map[string][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        body,
		Method:      "POST",
		Endpoint:    endpoint,
		QueryString: "",
		Headers:     headers,
	}
	d, ok := c.(*DeviceClient)
	if ok && d.IsMTLS {
		if err := req.setupForMTLS(d); err != nil {
			return nil, fmt.Errorf("Error setting up MTLS: %v", err)
		}
	}
	return do(c, req, creds)
}

func postWithCustomTransport(c cbClient, endpoint string, body interface{}, creds [][]string, headers map[string][]string, tr *http.Transport, isMTLS bool) (*CbResp, error) {
	req := &CbReq{
		Body:        body,
		Method:      "POST",
		Endpoint:    endpoint,
		QueryString: "",
		Headers:     headers,
		Transport:   tr,
		IsMTLS:      isMTLS,
	}
	return do(c, req, creds)
}

func put(c cbClient, endpoint string, body interface{}, heads [][]string, headers map[string][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        body,
		Method:      "PUT",
		Endpoint:    endpoint,
		QueryString: "",
		Headers:     headers,
	}
	d, ok := c.(*DeviceClient)
	if ok && d.IsMTLS {
		if err := req.setupForMTLS(d); err != nil {
			return nil, fmt.Errorf("Error setting up MTLS: %v", err)
		}
	}
	return do(c, req, heads)
}

func delete(c cbClient, endpoint string, query map[string]string, heads [][]string, headers map[string][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        nil,
		Method:      "DELETE",
		Endpoint:    endpoint,
		Headers:     headers,
		QueryString: query_to_string(query),
	}
	d, ok := c.(*DeviceClient)
	if ok && d.IsMTLS {
		if err := req.setupForMTLS(d); err != nil {
			return nil, fmt.Errorf("Error setting up MTLS: %v", err)
		}
	}
	return do(c, req, heads)
}

func deleteWithBody(c cbClient, endpoint string, body interface{}, heads [][]string, headers map[string][]string) (*CbResp, error) {
	req := &CbReq{
		Body:        body,
		Method:      "DELETE",
		Endpoint:    endpoint,
		Headers:     headers,
		QueryString: "",
	}
	d, ok := c.(*DeviceClient)
	if ok && d.IsMTLS {
		if err := req.setupForMTLS(d); err != nil {
			return nil, fmt.Errorf("Error setting up MTLS: %v", err)
		}
	}
	return do(c, req, heads)
}

func query_to_string(query map[string]string) string {
	qryVals := url.Values{}
	for k, v := range query {
		qryVals.Set(k, v)
	}
	return qryVals.Encode()
}

func checkForEdgeProxy(c cbClient, r *CbReq) {
	edgeProxy := c.getEdgeProxy()
	if r.Headers == nil {
		r.Headers = map[string][]string{}
	}
	if edgeProxy != nil {
		r.Headers["Clearblade-Systemkey"] = []string{edgeProxy.SystemKey}
		r.Headers["Clearblade-Edge"] = []string{edgeProxy.EdgeName}
	}
}

func parseEdgeConfig(e EdgeConfig) *exec.Cmd {
	cmd := exec.Command("edge",
		"-edge-ip=localhost",
		"-edge-id="+e.EdgeName,
		"-edge-cookie="+e.EdgeToken,
		"-platform-ip="+e.PlatformIP,
		"-platform-port="+e.PlatformPort,
		"-parent-system="+e.ParentSystem,
	)
	if p := e.HttpPort; p != "" {
		cmd.Args = append(cmd.Args, "-edge-listen-port="+p)
	}
	if p := e.MqttPort; p != "" {
		cmd.Args = append(cmd.Args, "-broker-tcp-port="+p)
	}
	if p := e.MqttTlsPort; p != "" {
		cmd.Args = append(cmd.Args, "-broker-tls-port="+p)
	}
	if p := e.WsPort; p != "" {
		cmd.Args = append(cmd.Args, "-broker-ws-port="+p)
	}
	if p := e.WssPort; p != "" {
		cmd.Args = append(cmd.Args, "-broker-wss-port="+p)
	}
	if p := e.AuthPort; p != "" {
		cmd.Args = append(cmd.Args, "-mqtt-auth-port="+p)
	}
	if p := e.AuthWsPort; p != "" {
		cmd.Args = append(cmd.Args, "-mqtt-ws-auth-port="+p)
	}
	if p := e.AdapterRootDir; p != "" {
		cmd.Args = append(cmd.Args, "-adaptors-root-dir="+p)
	}
	if e.Lean {
		cmd.Args = append(cmd.Args, "-lean-mode")
	}
	if e.Cache {
		cmd.Args = append(cmd.Args, "-local")
	}
	if p := e.LogLevel; p != "" {
		cmd.Args = append(cmd.Args, "-log-level="+p)
	}
	if e.Insecure {
		cmd.Args = append(cmd.Args, "-insecure=true")
	}
	if e.DevMode {
		cmd.Args = append(cmd.Args, "-development-mode=true")
	}
	if s := e.Stdout; s != nil {
		cmd.Stdout = s
	} else {
		cmd.Stdout = os.Stdout
	}
	if s := e.Stderr; s != nil {
		cmd.Stderr = s
	} else {
		cmd.Stderr = os.Stderr
	}
	return cmd
}

func makeSliceOfMaps(inIF interface{}) ([]map[string]interface{}, error) {
	switch inIF.(type) {
	case []interface{}:
		in := inIF.([]interface{})
		rval := make([]map[string]interface{}, len(in))
		for i, val := range in {
			valMap, ok := val.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("expected item to be a map, got %T", val)
			}
			rval[i] = valMap
		}
		return rval, nil
	case []map[string]interface{}:
		return inIF.([]map[string]interface{}), nil
	default:
		return nil, fmt.Errorf("Expected list of maps, got %T", inIF)
	}
}

func createQueryMap(query *Query) (map[string]string, error) {
	var qry map[string]string
	if query != nil {
		queryMap := query.serialize()
		queryBytes, err := json.Marshal(queryMap)
		if err != nil {
			return nil, err
		}
		qry = map[string]string{
			"query": string(queryBytes),
		}
	} else {
		qry = nil
	}

	return qry, nil
}

func nicelySetRefreshToken(body map[string]interface{}) string {
	if tok, ok := body["refresh_token"].(string); ok {
		return tok
	}
	return ""
}

func decodeMapToStruct(incoming interface{}, outgoing interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  outgoing,
		TagName: "json",
	})

	if err != nil {
		return err
	}

	return decoder.Decode(incoming)
}
