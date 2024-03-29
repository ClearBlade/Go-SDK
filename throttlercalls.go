package GoSDK

//"fmt"

const (
	_THROTTLERS_PREAMBLE    = "/admin/throttlers"
	_USERS_DEVICES_PREAMBLE = "/admin/usersanddevices"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) GetAllThrottlers() (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := mapResponse(get(d, _THROTTLERS_PREAMBLE, nil, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetThrottler(throttlerName string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := mapResponse(get(d, _THROTTLERS_PREAMBLE+"/"+throttlerName, nil, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) DeleteThrottlerCasesAndExceptions(throttlerName string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := mapResponse(delete(d, _THROTTLERS_PREAMBLE+"/"+throttlerName, nil, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) CreateThrottlerCase(throttlerName string, caseInfo map[string]interface{}) (map[string]interface{}, error) {
	return d.createThrottlerCase(throttlerName, "/cases", caseInfo)
}

func (d *DevClient) CreateThrottlerException(throttlerName string, caseInfo map[string]interface{}) (map[string]interface{}, error) {
	return d.createThrottlerCase(throttlerName, "/exceptions", caseInfo)
}

func (d *DevClient) createThrottlerCase(throttlerName, pathTail string, caseInfo map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := mapResponse(post(d, _THROTTLERS_PREAMBLE+"/"+throttlerName+pathTail, caseInfo, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeleteAllThrottlerCases(throttlerName string) (map[string]interface{}, error) {
	return d.deleteAllThrottlerCases(throttlerName, "/cases")
}

func (d *DevClient) DeleteAllThrottlerExceptions(throttlerName string) (map[string]interface{}, error) {
	return d.deleteAllThrottlerCases(throttlerName, "/exceptions")
}

func (d *DevClient) deleteAllThrottlerCases(throttlerName, pathTail string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := mapResponse(delete(d, _THROTTLERS_PREAMBLE+"/"+throttlerName+pathTail, nil, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) UpdateThrottlerCase(throttlerName, caseName string, updates map[string]interface{}) (map[string]interface{}, error) {
	return d.updateThrottlerCase(throttlerName, caseName, "/cases", updates)
}

func (d *DevClient) UpdateThrottlerException(throttlerName, caseName string, updates map[string]interface{}) (map[string]interface{}, error) {
	return d.updateThrottlerCase(throttlerName, caseName, "/exceptions", updates)
}

func (d *DevClient) updateThrottlerCase(throttlerName, caseName, pathTail string, updates map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := mapResponse(put(d, _THROTTLERS_PREAMBLE+"/"+throttlerName+pathTail+"/"+caseName, updates, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeleteThrottlerCase(throttlerName, caseName string) (map[string]interface{}, error) {
	return d.deleteThrottlerCase(throttlerName, caseName, "/cases")
}

func (d *DevClient) DeleteThrottlerException(throttlerName, caseName string) (map[string]interface{}, error) {
	return d.deleteThrottlerCase(throttlerName, caseName, "/exceptions")
}

func (d *DevClient) deleteThrottlerCase(throttlerName, caseName, pathTail string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := mapResponse(delete(d, _THROTTLERS_PREAMBLE+"/"+throttlerName+pathTail+"/"+caseName, nil, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (d *DevClient) AllUsersAndDevices(inputText string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := mapResponse(get(d, _USERS_DEVICES_PREAMBLE+"?text="+inputText, nil, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) OneSystemUsersAndDevices(systemKey, inputText string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := mapResponse(get(d, _USERS_DEVICES_PREAMBLE+"/"+systemKey+"?text="+inputText, nil, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}
