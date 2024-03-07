package GoSDK

import "fmt"

const (
	_SYSTEM_MIGRATION_DEV_PREAMBLE = "/api/v/1/systemmigration/"
)

func (d *DevClient) UploadToSystem(systemKey string, zipBuffer []byte, dryRun bool) (interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	zipContents, err := getContentsForFile(zipBuffer)
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{
		"contents": zipContents,
	}

	url := fmt.Sprintf("%s%s/upload?dryRun=%t", _SYSTEM_MIGRATION_DEV_PREAMBLE, systemKey, dryRun)
	resp, err := post(d, url, body, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (d *DevClient) GetSystemUploadVersion(systemKey string) (interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s/upload", _SYSTEM_MIGRATION_DEV_PREAMBLE, systemKey)
	resp, err := get(d, url, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
