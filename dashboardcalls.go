package GoSDK

const (
	_DASHBOARDS_DEV_PREAMBLE  = "/admin/dashboards/"
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
