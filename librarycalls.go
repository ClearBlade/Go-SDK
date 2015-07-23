package GoSDK

import (
	//"fmt"
	"strconv"
)

const (
	_LIB_PREAMBLE  = "/codeadmin/v/2/library"
	_HIST_PREAMBLE = "/codeadmin/v/2/history/library"
)

func (d *DevClient) GetAllLibraries() ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_LIB_PREAMBLE, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetLibraries(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_LIB_PREAMBLE+"/"+systemKey, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetLibrary(systemKey, name string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_LIB_PREAMBLE+"/"+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) CreateLibrary(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(_LIB_PREAMBLE+"/"+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) UpdateLibrary(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(_LIB_PREAMBLE+"/"+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeleteLibrary(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(_LIB_PREAMBLE+"/"+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) GetVersionHistory(systemKey, name string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_HIST_PREAMBLE+"/"+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetVersion(systemKey, name string, version int) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_HIST_PREAMBLE+"/"+systemKey+"/"+name+"/"+strconv.Itoa(version), nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}
