package GoSDK

func (d *DevClient) EnterProvisioningMode(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := post(d, "/admin/"+systemKey+"/edgemode/provisioning", map[string]interface{}{}, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}

	return resp.Body.([]interface{}), nil
}

func (d *DevClient) EnterRuntimeMode(systemKey string, info map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := post(d, "/admin/"+systemKey+"/edgemode/runtime", info, creds, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}
