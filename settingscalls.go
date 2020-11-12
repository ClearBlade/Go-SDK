package GoSDK

import (
//"fmt"
)

const (
	_SETTINGS_PREAMBLE = "/admin/settings/"
	_GOOGLE_SERVICE    = "google-storage"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) GetGoogleStorageSettings(systemKey string) (map[string]interface{}, error) {
	return getGoogleStorageSettings(d, _SETTINGS_PREAMBLE+_GOOGLE_SERVICE)
}

func getGoogleStorageSettings(c cbClient, endpoint string) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(c, endpoint, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) AddGoogleStorageSettings(settings map[string]interface{}) (map[string]interface{}, error) {
	return addGoogleStorageSettings(d, _SETTINGS_PREAMBLE+_GOOGLE_SERVICE, settings)
}

func addGoogleStorageSettings(c cbClient, endpoint string, settings map[string]interface{}) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(c, endpoint, settings, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) UpdateGoogleStorageSettings(settings map[string]interface{}) (map[string]interface{}, error) {
	return updateGoogleStorageSettings(d, _SETTINGS_PREAMBLE+_GOOGLE_SERVICE, settings)
}

func updateGoogleStorageSettings(c cbClient, endpoint string, settings map[string]interface{}) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(c, endpoint, settings, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) DeleteGoogleStorageSettings(systemKey, deploymentName, box, relPath string) error {
	return deleteGoogleStorageSettings(d, _SETTINGS_PREAMBLE+_GOOGLE_SERVICE)
}

func deleteGoogleStorageSettings(c cbClient, endpoint string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	_, err = delete(c, endpoint, nil, creds, nil)
	return err
}
