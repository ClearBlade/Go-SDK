package GoSDK

func (d *DevClient) GetEdgeGroups(systemKey string, query *Query) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	qry, err := createQueryMap(query)
	if err != nil {
		return nil, err
	}

	resp, err := get(d, "/admin/"+systemKey+"/edge_groups", qry, creds, nil)
	resp, err = mapResponse(resp, err)

	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetEdgeGroup(systemKey, name string, recursive bool) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	var queries map[string]string
	if recursive == true {
		queries = map[string]string{"recursive": "true"}
	}
	resp, err := get(d, "/admin/"+systemKey+"/edge_groups/"+name, queries, creds, nil)
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
