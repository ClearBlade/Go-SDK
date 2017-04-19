package GoSDK

import (
	"encoding/base64"
	"fmt"
)

const (
	_ADAPTOR_HEADER_KEY    = "ClearBlade-DeviceToken"
	_ADAPTORS_DEV_PREAMBLE = "/admin/"
)

func (d *DevClient) GetAdaptors(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _ADAPTORS_DEV_PREAMBLE+systemKey+"/adaptors", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetAdaptor(systemKey, name string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _ADAPTORS_DEV_PREAMBLE+systemKey+"/adaptors/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) CreateAdaptor(systemKey, name string,
	data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	data["name"] = name
	resp, err := post(d, _ADAPTORS_DEV_PREAMBLE+systemKey+"/adaptors", data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeleteAdaptor(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, _ADAPTORS_DEV_PREAMBLE+systemKey+"/adaptors/"+name, nil, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

func (d *DevClient) UpdateAdaptor(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(d, _ADAPTORS_DEV_PREAMBLE+systemKey+"/adaptors/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetAdaptorFiles(systemKey, adaptorName string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _ADAPTORS_DEV_PREAMBLE+systemKey+"/adaptors/"+adaptorName+"/files", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetAdaptorFile(systemKey, adaptorName, fileName string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _ADAPTORS_DEV_PREAMBLE+systemKey+"/adaptors/"+adaptorName+"/files/"+fileName, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) CreateAdaptorFile(systemKey, adaptorName, fileName string,
	data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	// the file member of 'data' could (and probably is) a binary executable file or library. Need to encode/decode this
	// for transfer
	contentsIF, exists := data["file"]
	if !exists {
		return nil, fmt.Errorf("'file' key/value pair missing in CreateAdaptorFile")
	}
	byts, ok := contentsIF.([]byte)
	if !ok {
		return nil, fmt.Errorf("Bad type for 'file' k/v pair in CreateAdaptorFile: %T: (%+v)", contentsIF, contentsIF)
	}
	data["file"] = base64.StdEncoding.EncodeToString(byts)
	data["name"] = fileName
	data["adaptor_name"] = adaptorName

	resp, err := post(d, _ADAPTORS_DEV_PREAMBLE+systemKey+"/adaptors/"+adaptorName+"/files", data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeleteAdaptorFile(systemKey, adaptorName, fileName string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, _ADAPTORS_DEV_PREAMBLE+systemKey+"/adaptors/"+adaptorName+"/files/"+fileName, nil, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

func (d *DevClient) UpdateAdaptorFile(systemKey, adaptorName, fileName string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(d, _ADAPTORS_DEV_PREAMBLE+systemKey+"/adaptors/"+adaptorName+"/files/"+fileName, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeployAdaptor(systemKey, adaptorName string, deploySpec map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(d, _ADAPTORS_DEV_PREAMBLE+systemKey+"/adaptors/"+adaptorName+"/deploy", deploySpec, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) ControlAdaptor(systemKey, adaptorName string, controlSpec map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(d, _ADAPTORS_DEV_PREAMBLE+systemKey+"/adaptors/"+adaptorName+"/control", controlSpec, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}
