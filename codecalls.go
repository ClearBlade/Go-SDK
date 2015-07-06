package GoSDK

import (
	"fmt"
)

const (
	_CODE_PREAMBLE = "/api/v/1/code"
)

func callService(c cbClient, systemKey, name string, params map[string]interface{}, log bool) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	var resp *CbResp
	if log {

		resp, err = post(_CODE_PREAMBLE+"/"+systemKey+"/"+name, params, creds, map[string][]string{"Logging_enabled": []string{"true"}})
	} else {
		resp, err = post(_CODE_PREAMBLE+"/"+systemKey+"/"+name, params, creds, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("Error calling %s service: %v", name, err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error calling %s service: %v", name, resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) CallService(systemKey, name string, params map[string]interface{}, log bool) (map[string]interface{}, error) {
	return callService(d, systemKey, name, params, log)
}

func (u *UserClient) CallService(systemKey, name string, params map[string]interface{}) (map[string]interface{}, error) {
	return callService(u, systemKey, name, params, false)
}
