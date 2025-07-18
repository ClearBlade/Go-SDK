package GoSDK

import (
	"encoding/json"
	"fmt"
)

//"fmt"

const (
	_SETTINGS_PREAMBLE = "/admin/settings/"
	_GOOGLE_SERVICE    = "google-storage"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) GetGoogleStorageSettings() (map[string]interface{}, error) {
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

func (d *DevClient) DeleteGoogleStorageSettings() error {
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) AddMTLSSettings(rootCA, crl string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	settings := map[string]interface{}{
		"root_ca": rootCA,
		"crl":     crl,
	}
	resp, err := put(d, _SETTINGS_PREAMBLE+"mtls", settings, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) AddMTLSSystemCertificate(systemKey, name, systemCA, privateKey string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	settings := map[string]interface{}{
		"system_key":  systemKey,
		"certificate": systemCA,
		"name":        name,
	}
	if privateKey != "" {
		settings["private_key"] = privateKey
	}
	resp, err := post(d, _SETTINGS_PREAMBLE+"mtls/"+systemKey+"/"+name, settings, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) GenerateRootCACertificate(keyType string, keySize int, expiryTimestamp int64) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"key_type":         keyType,
		"key_size":         keySize,
		"expiry_timestamp": expiryTimestamp,
	}
	resp, err := post(d, _SETTINGS_PREAMBLE+"mtls/generate-certificate", payload, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) GenerateSystemCertificate(systemKey, name, keyType string, keySize int, expiryTimestamp int64) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"key_type":         keyType,
		"key_size":         keySize,
		"expiry_timestamp": expiryTimestamp,
	}
	resp, err := post(d, _SETTINGS_PREAMBLE+"mtls/"+systemKey+"/"+name+"/generate-certificate", payload, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) UpdateMTLSSystemCertificate(systemKey, name, systemCA, privateKey string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	settings := map[string]interface{}{
		"system_key":  systemKey,
		"certificate": systemCA,
		"name":        name,
	}
	if privateKey != "" {
		settings["private_key"] = privateKey
	}
	resp, err := put(d, _SETTINGS_PREAMBLE+"mtls/"+systemKey+"/"+name, settings, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) GetMTLSSystemCertificate(systemKey, name string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _SETTINGS_PREAMBLE+"mtls/"+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetAllMTLSSystemCertificates(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _SETTINGS_PREAMBLE+"mtls/"+systemKey, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) DeleteMTLSSystemCertificate(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, _SETTINGS_PREAMBLE+"mtls/"+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) GetMTLSSettings() (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _SETTINGS_PREAMBLE+"mtls", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeleteMTLSSettings() error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, _SETTINGS_PREAMBLE+"mtls", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) DeleteRevokedCertificate(certificateHash string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	payload := map[string]string{
		"certificate_hash": certificateHash,
	}
	resp, err := delete(d, "/admin/revoked_certs", payload, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) RevokeCertificate(certificateHash string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	payload := map[string]interface{}{
		"certificate_hash": certificateHash,
	}
	resp, err := post(d, "/admin/revoked_certs", payload, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}

func (d *DevClient) GetRevokedCertificates(query *Query) ([]map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	var qry map[string]string
	if query != nil {
		query_map := query.serialize()
		query_bytes, err := json.Marshal(query_map)
		if err != nil {
			return nil, err
		}
		qry = map[string]string{
			"query": string(query_bytes),
		}
	} else {
		qry = nil
	}
	resp, err := get(d, "/admin/revoked_certs", qry, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting data: %v", err)
	}
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]map[string]interface{}), nil
}

func (d *DevClient) DeleteRevokedCertificates(query *Query) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	var qry map[string]string
	if query != nil {
		query_map := query.serialize()
		query_bytes, err := json.Marshal(query_map)
		if err != nil {
			return err
		}
		qry = map[string]string{
			"query": string(query_bytes),
		}
	} else {
		qry = nil
	}
	resp, err := delete(d, "/admin/revoked_certs", qry, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}
