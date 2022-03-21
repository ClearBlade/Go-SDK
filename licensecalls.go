package GoSDK

func (d *DevClient) UpdateLicense(changes map[string]interface{}) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	body := map[string]interface{}{"license": changes}

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
