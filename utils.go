package GoSDK

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	addr string
)

type CbCli struct {
	URL     string
	Headers map[string]string
}

type CbReq struct {
	Body        interface{}
	Method      string
	Endpoint    string
	QueryString string
}

type CbResp struct {
	Body       interface{}
	StatusCode int
}

func NewCbCli() *CbCli {
	return &CbCli{
		URL:     "https://platform.clearblade.com",
		Headers: map[string]string{},
	}
}

func (c *CbCli) AddHeader(key, value string) {
	c.Headers[key] = value
}

func (c *CbCli) RemoveHeader(key string) {
	delete(c.Headers, key)
}

func (c *CbCli) SetSystem(key, secret string) {
	c.AddHeader("SystemKey", key)
	c.AddHeader("SystemSecret", secret)
}

func (c *CbCli) SetDevToken(tok string) {
	c.RemoveHeader("SystemKey")
	c.RemoveHeader("SystemSecret")
	c.RemoveHeader("ClearBlade-UserToken") // just in case
	c.AddHeader("ClearBlade-DevToken", tok)
}

func (c *CbCli) SetUserToken(tok string) {
	c.RemoveHeader("SystemKey")
	c.RemoveHeader("SystemSecret")
	c.RemoveHeader("ClearBlade-DevToken") // just in case
	c.AddHeader("ClearBlade-UserToken", tok)
}

func (c *CbCli) Do(r *CbReq) (*CbResp, error) {
	var bodyToSend *bytes.Buffer
	if r.Body != nil {
		b, jsonErr := json.Marshal(r.Body)
		if jsonErr != nil {
			return nil, fmt.Errorf("JSON Encoding Error: %v\n", jsonErr)
		}
		bodyToSend = bytes.NewBuffer(b)
	} else {
		bodyToSend = nil
	}
	url := c.URL + r.Endpoint
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
		return nil, fmt.Errorf("Request Creation Error: %v\n", reqErr)
	}
	for k, v := range c.Headers {
		req.Header.Add(k, v)
	}

	cli := &http.Client{}
	log.Printf("CLI: %+v\n", cli)
	log.Printf("REQ: %+v\n", req)
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error Making Request: %v\n", err)
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("Error Reading Response Body: %v\n", readErr)
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

func (c *CbCli) Get(endpoint string, query map[string]string) (*CbResp, error) {
	req := &CbReq{
		Body:        nil,
		Method:      "GET",
		Endpoint:    endpoint,
		QueryString: query_to_string(query),
	}
	return c.Do(req)
}

func (c *CbCli) Post(endpoint string, body interface{}) (*CbResp, error) {
	req := &CbReq{
		Body:        body,
		Method:      "POST",
		Endpoint:    endpoint,
		QueryString: "",
	}
	return c.Do(req)
}

func (c *CbCli) Put(endpoint string, body interface{}) (*CbResp, error) {
	req := &CbReq{
		Body:        body,
		Method:      "PUT",
		Endpoint:    endpoint,
		QueryString: "",
	}
	return c.Do(req)
}

func (c *CbCli) Delete(endpoint string, query map[string]string) (*CbResp, error) {
	req := &CbReq{
		Body:        nil,
		Method:      "DELETE",
		Endpoint:    endpoint,
		QueryString: query_to_string(query),
	}
	return c.Do(req)
}

func query_to_string(query map[string]string) string {
	qryStr := ""
	for k, v := range query {
		qryStr += k + "=" + v + "&"
	}
	return strings.TrimSuffix(qryStr, "&")
}
