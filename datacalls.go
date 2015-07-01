package GoSDK

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	_DATA_PREAMBLE      = "/api/v/1/data/"
	_DATA_NAME_PREAMBLE = "/api/v/1/collection/"
)

func (u *UserClient) InsertData(collection_id string, data interface{}) error {
	return insertdata(u, collection_id, data)
}

func (d *DevClient) InsertData(collection_id string, data interface{}) error {
	return insertdata(d, collection_id, data)
}

func insertdata(c cbClient, collection_id string, data interface{}) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	resp, err := post(_DATA_PREAMBLE+collection_id, data, creds, nil)
	if err != nil {
		return fmt.Errorf("Error inserting: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error inserting: %v", resp.Body)
	}
	return nil
}

func (u *UserClient) GetData(collection_id string, query *Query) (map[string]interface{}, error) {
	return getdata(u, collection_id, query)
}

func (u *UserClient) GetDataByName(collectionName string, query *Query) (map[string]interface{}, error) {
	return getDataByName(u, u.SystemKey, collectionName, query)
}

func (d *DevClient) GetDataByName(collectionName string, query *Query) (map[string]interface{}, error) {
	return nil, fmt.Errorf("Developer cannot call this yet")
}

func (d *DevClient) GetData(collection_id string, query *Query) (map[string]interface{}, error) {
	return getdata(d, collection_id, query)
}

func getDataByName(c cbClient, sysKey string, collectionName string, query *Query) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	var qry map[string]string
	if query != nil {
		query_map := query.serialize()
		query_bytes, err := json.Marshal(query_map)
		if err != nil {
			return nil, err
		}
		qry = map[string]string{
			"query": url.QueryEscape(string(query_bytes)),
		}
	} else {
		qry = nil
	}
	resp, err := get(_DATA_NAME_PREAMBLE+sysKey+"/"+collectionName, qry, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting data: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting data: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func getdata(c cbClient, collection_id string, query *Query) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	var qry map[string]string
	if query != nil {
		query_map := query.serialize()
		query_bytes, err := json.Marshal(query_map)
		if err != nil {
			return nil, err
		}
		qry = map[string]string{
			"query": url.QueryEscape(string(query_bytes)),
		}
	} else {
		qry = nil
	}
	resp, err := get(_DATA_PREAMBLE+collection_id, qry, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting data: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting data: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (u *UserClient) UpdateData(collection_id string, query *Query, changes map[string]interface{}) error {
	err := updatedata(u, collection_id, query, changes)
	return err
}

func (d *DevClient) UpdateData(collection_id string, query *Query, changes map[string]interface{}) error {
	err := updatedata(d, collection_id, query, changes)
	return err
}

func updatedata(c cbClient, collection_id string, query *Query, changes map[string]interface{}) error {
	qry := query.serialize()
	body := map[string]interface{}{
		"query": qry,
		"$set":  changes,
	}
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	resp, err := put(_DATA_PREAMBLE+collection_id, body, creds, nil)
	if err != nil {
		return fmt.Errorf("Error updating data: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating data: %v", resp.Body)
	}
	return nil
}

func (u *UserClient) DeleteData(collection_id string, query *Query) error {
	return deletedata(u, collection_id, query)
}

func (d *DevClient) DeleteData(collection_id string, query *Query) error {
	return deletedata(d, collection_id, query)
}

func deletedata(c cbClient, collection_id string, query *Query) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	var qry map[string]string
	if query != nil {
		query_map := query.serialize()
		query_bytes, err := json.Marshal(query_map)
		if err != nil {
			return err
		}
		qry = map[string]string{
			"query": url.QueryEscape(string(query_bytes)),
		}
	} else {
		qry = nil
	}
	resp, err := delete(_DATA_PREAMBLE+collection_id, qry, creds, nil)
	if err != nil {
		return fmt.Errorf("Error deleting data: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting data: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) GetColumns(collection_id string) ([]interface{}, error) {
	return getColumns(d, collection_id)
}

func (u *UserClient) GetColumns(collection_id string) ([]interface{}, error) {
	return getColumns(u, collection_id)
}

func getColumns(c cbClient, collection_id string) ([]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(_DATA_PREAMBLE+collection_id+"/columns", nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting collection columns: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting collection columns: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetDataByKeyAndName(string, string, *Query) (map[string]interface{}, error) {
	return nil, fmt.Errorf("Unimplemented")
}

func (d *UserClient) GetDataByKeyAndName(string, string, *Query) (map[string]interface{}, error) {
	return nil, fmt.Errorf("Unimplemented")

}
