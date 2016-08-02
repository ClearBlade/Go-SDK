package GoSDK

const (
	_DASHBOARDS_DEV_PREAMBLE  = "/admin/portals/"
	_DASHBOARDS_USER_PREAMBLE = "/api/v/2/portals/"
)

func (d *DevClient) GetDashboards(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _DASHBOARDS_DEV_PREAMBLE+systemKey, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetDashboard(systemKey, name string) (map[string]interface{}, error){
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _DASHBOARDS_USER_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) CreateDashboard(systemKey, name string, dash map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := post(d, _DASHBOARDS_USER_PREAMBLE+systemKey+"/"+name, make(map[string]interface{}), creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	// An empty dashbaord has now been created.. populate with provided data
	return d.UpdateDashboard(systemKey, name, dash)
}

func (d *DevClient) UpdateDashboard(systemKey, name string, dash map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(d, _DASHBOARDS_USER_PREAMBLE+systemKey+"/"+name, dash, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeleteDashboard(systemKey, name string) (map[string]interface{}, error){
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := delete(d, _DASHBOARDS_DEV_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}
