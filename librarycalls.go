package GoSDK

import (
	"fmt"
)

const (
	_LIB_PREAMBLE = "/codeadmin/v/2/library/"
)

func (d *DevClient) GetLibraries(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_LIB_PREAMBLE+systemKey, nil, creds)

	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	fmt.Printf("GOT BACK: %+v\n", resp.Body)
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) CreateLibrary(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(_LIB_PREAMBLE+systemKey+"/"+name, data, creds)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}

	return resp.Body.(map[string]interface{}), nil
}
