package GoSDK

import "fmt"

const (
	_ABOUT_PREAMBLE = "/api/v/1/about"
)

func (d *DevClient) GetConfiguration() (interface{}, error) {
	return d.getAboutEndpoint("configuration")
}

func (d *DevClient) GetDatabaseStats(systemKey string) (interface{}, error) {
	return d.getAboutEndpoint("dbstats")
}

func (d *DevClient) getAboutEndpoint(endpoint string) (interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%s", _ABOUT_PREAMBLE, endpoint)
	resp, err := get(d, url, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
