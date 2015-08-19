package GoSDK

import (
	//"fmt"
	"strconv"
)

const (
	_LIB_PREAMBLE  = "/codeadmin/v/2/library"
	_HIST_PREAMBLE = "/codeadmin/v/2/history/library"
)

//GetAllLibraries returns descriptions of all libraries for a platform instance
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

//GetLibrary returns a list of libraries for a system
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

//GetLibrary returns information about a specific library
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

//CreateLibrary allows the developer to create a library to be called by other service functions
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

//UpdateLibrary allows the developer to change the content of the library
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

//DeleteLibrary allows the developer to remove library content
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

//GetVersionHistory returns a slice of library descriptions corresponding to each library
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

//GetVersion gets the current version of a library
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
