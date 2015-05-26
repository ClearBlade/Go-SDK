package GoSDK

import (
	"fmt"
)

const (
	_EVENTS_DEFS_PREAMBLE  = "/api/v/1/events/definitions"
	_EVENTS_HDLRS_PREAMBLE = "/api/v/1/events/handlers/"
)

func (d *DevClient) GetEventDefinitions() ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_EVENTS_DEFS_PREAMBLE, nil, creds)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	fmt.Printf("IT WORKED: %+v\n", resp)
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetEventHandlers(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_EVENTS_HDLRS_PREAMBLE+systemKey, nil, creds)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) CreateEventHandler(systemKey string,
	data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(_EVENTS_HDLRS_PREAMBLE+systemKey, data, creds)
	resp, err = mapResponse(resp, err)
	fmt.Printf("RESP: %+v\n", resp)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func mapResponse(resp *CbResp, err error) (*CbResp, error) {
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(resp.Body.(string))
	}
	return resp, nil
}
