package GoSDK

const (
	_DASHBOARDS_DEV_PREAMBLE  = "/admin/dashboards/"
	_DASHBOARDS_USER_PREAMBLE = "/api/v/2/dashboards/"
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

func (d *DevClient) CreateDashboard(systemKey, name string, dash map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(d, _DASHBOARDS_USER_PREAMBLE+systemKey+"/"+name, dash, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}


