package GoSDK

import (
	"fmt"
)

const (
	_SIDECAR_PREAMBLE_ = "/ia-sidecar"
)

func (d *DevClient) CreateIAInstance(name string, data map[string]interface{}) (map[string]interface{}, error) {
	return createIAInstance(d, _SIDECAR_PREAMBLE_+"/instances", data)
}

func createIAInstance(c cbClient, endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
	fmt.Printf("createIAInstance\n")
	creds, err := c.credentials()
	fmt.Printf("\ncreds: %+v\n\nerr:%+v\n", creds, err)
	if err != nil {
		return nil, err
	}
	resp, err := post(c, endpoint, data, creds, nil)
	fmt.Printf("resp: %+v, err: %+v\n", resp, err)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}
