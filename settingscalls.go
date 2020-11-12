package GoSDK

import (
//"fmt"
)

const (
	_SETTINGS_PREAMBLE = "/admin/settings/"
	_GOOGLE_SERVICE    = "google-storage"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) GetGoogleServiceSettings(systemKey string) (map[string]interface{}, error) {
	return getGoogleServiceSettings(d, _SETTINGS_PREAMBLE+_GOOGLE_SERVICE)
}

func getGoogleServiceSettings(c cbClient, endpoint string) (map[string]interface{}, error) {
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

func (d *DevClient) AddGoogleServiceSettings(settings map[string]interface{}) (map[string]interface{}, error) {
	return addGoogleServiceSettings(d, _SETTINGS_PREAMBLE+_GOOGLE_SERVICE, settings)
}

func addGoogleServiceSettings(c cbClient, endpoint string, settings map[string]interface{}) (map[string]interface{}, error) {
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

func (d *DevClient) UpdateGoogleServiceSettings(settings map[string]interface{}) (map[string]interface{}, error) {
	return updateGoogleServiceSettings(d, _SETTINGS_PREAMBLE+_GOOGLE_SERVICE, settings)
}

func updateGoogleServiceSettings(c cbClient, endpoint string, settings map[string]interface{}) (map[string]interface{}, error) {
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

func (d *DevClient) DeleteGoogleServiceSettings(systemKey, deploymentName, box, relPath string) error {
	return deleteGoogleServiceSettings(d, _SETTINGS_PREAMBLE+_GOOGLE_SERVICE)
}

func deleteGoogleServiceSettings(c cbClient, endpoint string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	_, err = delete(c, endpoint, nil, creds, nil)
	return err
}
