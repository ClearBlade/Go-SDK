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
