package GoSDK

const (
	_SIDECAR_PREAMBLE_ = "/ia-sidecar"
)

func (d *DevClient) CreateIAInstance(name string, data map[string]interface{}) (map[string]interface{}, error) {
	return createIAInstance(d, _SIDECAR_PREAMBLE_+"/instances", data)
}

func createIAInstance(c cbClient, endpoint string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(c, endpoint, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}
