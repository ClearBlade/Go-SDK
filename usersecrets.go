package GoSDK

const (
	_USER_SECRETS_PREAMBLE = "/api/v/4/user_secrets/"
)

func (d *DevClient) AddSecret(systemKey, name string, data interface{}) (string, error) {
	creds, err := d.credentials()
	if err != nil {
		return "", err
	}
	resp, err := post(d, _USER_SECRETS_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return "", err
	}
	return resp.Body.(string), nil
}

func (d *DevClient) UpdateSecret(systemKey, name string, data interface{}) (string, error) {
	creds, err := d.credentials()
	if err != nil {
		return "", err
	}
	resp, err := put(d, _USER_SECRETS_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return "", err
	}
	return resp.Body.(string), nil
}

func (d *DevClient) GetSecrets(systemKey string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _USER_SECRETS_PREAMBLE+systemKey, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetSecret(systemKey, name string) (interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return "", err
	}
	resp, err := get(d, _USER_SECRETS_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return "", err
	}
	return resp.Body, nil
}

func (d *DevClient) DeleteSecret(systemKey, name string) (string, error) {
	creds, err := d.credentials()
	if err != nil {
		return "", err
	}
	resp, err := delete(d, _USER_SECRETS_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return "", err
	}
	return resp.Body.(string), nil
}

func (d *DevClient) DeleteSecrets(systemKey string) (string, error) {
	creds, err := d.credentials()
	if err != nil {
		return "", err
	}
	resp, err := delete(d, _USER_SECRETS_PREAMBLE+systemKey, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return "", err
	}
	return resp.Body.(string), nil
}
