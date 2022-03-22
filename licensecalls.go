package GoSDK

func (d *DevClient) UpdateLicense(changes string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	body := map[string]interface{}{"pkey": changes}

	resp, err := put(d, "/admin/pkey", body, creds, nil)
	_, err = mapResponse(resp, err)
	if err != nil {
		return err
	}

	return nil
}

func (d *DevClient) ConnectedNodes() ([]string, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(d, "/admin/connectednodes", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}

	return resp.Body.([]string), nil
}

func (d *DevClient) ConnectedEdges() ([]map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(d, "/admin/connectededges", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}

	return resp.Body.([]map[string]interface{}), nil
}

func (d *DevClient) LicenseLimits() (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	limits, err := get(d, "/admin/limits", nil, creds, nil)
	limits, err = mapResponse(limits, err)
	if err != nil {
		return nil, err
	}
	return limits.Body.(map[string]interface{}), nil
}
