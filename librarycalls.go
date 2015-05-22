package GoSDK

import (
	"fmt"
)

const (
	_LIB_PREAMBLE = "/codeadmin/api/v/1/library/"
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
