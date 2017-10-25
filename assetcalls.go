package GoSDK

func (d *DevClient) GetSystemAssetDeployments(systemKey string, query *Query) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	qry, err := createQueryMap(query)
	if err != nil {
		return nil, err
	}

	resp, err := get(d, "/admin/"+systemKey+"/deploy_assets", qry, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}

	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetAssetClassDeployments(systemKey, assetClass string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(d, "/admin/"+systemKey+"/deploy_assets/"+assetClass, nil, creds, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) UpdateAssetClassDeployments(systemKey, assetClass string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := put(d, "/admin/"+systemKey+"/deploy_assets/"+assetClass, data, creds, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetAssetDeployments(systemKey, assetClass, assetId string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(d, "/admin/"+systemKey+"/deploy_assets/"+assetClass+"/"+assetId, nil, creds, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) UpdateAssetDeployments(systemKey, assetClass, assetId string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := put(d, "/admin/"+systemKey+"/deploy_assets/"+assetClass+"/"+assetId, data, creds, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetAssetsDeployedToEntity(systemKey, entityType, entityName string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	headers := map[string][]string{entityType: []string{entityName}}

	resp, err := put(d, "/admin/"+systemKey+"/deployed_assets", nil, creds, headers)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetAssetsNotDeployedOnPlatform(systemKey string, query *Query) ([]interface{}, error) {

	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	qry, err := createQueryMap(query)
	if err != nil {
		return nil, err
	}

	resp, err := get(d, "/admin/"+systemKey+"/deploy_on_platform", qry, creds, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetAssetPlatformDeploymentStatus(systemKey, assetClass, assetId string) (map[string]interface{}, error) {

	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(d, "/admin/"+systemKey+"/deploy_on_platform/"+assetClass+"/"+assetId, nil, creds, nil)
	if err != nil {
		return nil, err
	}

	return resp.Body.(map[string]interface{}), err
}

func (d *DevClient) UpdateAssetPlatformDeploymentStatus(systemKey, assetClass, assetId string, deploy bool) (map[string]interface{}, error) {

	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	reqBody := map[string]interface{}{"deploy": deploy}
	resp, err := put(d, "/admin/"+systemKey+"/deploy_on_platform/"+assetClass+"/"+assetId, reqBody, creds, nil)
	if err != nil {
		return nil, err
	}

	return resp.Body.(map[string]interface{}), nil
}

/* just using this as boilerplate example

func (d *DevClient) GetEdgeGroup(systemKey, name string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, "/admin/"+systemKey+"/edge_groups/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) CreateEdgeGroup(systemKey, name string,
	data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	data["name"] = name
	resp, err := post(d, "/admin/"+systemKey+"/edge_groups", data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) DeleteEdgeGroup(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, "/admin/"+systemKey+"/edge_groups/"+name, nil, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

func (d *DevClient) UpdateEdgeGroup(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(d, "/admin/"+systemKey+"/edge_groups/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func mkE1() map[string]interface{} {
	return map[string]interface{}{
		"system_key":    "",
		"system_secret": "Fred's system secret",
		"token":         "Fred's Token",
		"description":   "Fred makes widgets",
		"location":      "123 4th Street",
	}
}

func mkE2() map[string]interface{} {
	return map[string]interface{}{
		"system_key":    "Pat's system key",
		"system_secret": "Pat's system secret",
		"token":         "Pat's Token",
		"description":   "Pat builds a biology lab",
		"location":      "999 Caesar Chavez",
	}
}

func mkE3() map[string]interface{} {
	return map[string]interface{}{
		"system_key":    "Adam's system key",
		"system_secret": "Adam's system secret",
		"token":         "Adam's Token",
		"description":   "Adam creates design patterns",
		"location":      "Somewhere on I35",
	}
}
*/
