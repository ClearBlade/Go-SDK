package GoSDK

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	_DATA_PREAMBLE = "/api/v/1/data/"
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
	resp, err := post(_DATA_PREAMBLE+collection_id, data, creds)
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

func (d *DevClient) GetData(collection_id string, query *Query) (map[string]interface{}, error) {
	return getdata(d, collection_id, query)
}

func getdata(c cbClient, collection_id string, query *Query) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	query_map := query.serialize()
	query_bytes, err := json.Marshal(query_map)
	if err != nil {
		return nil, err
	}
	qry := map[string]string{
		"query": url.QueryEscape(string(query_bytes)),
	}
	resp, err := get(_DATA_PREAMBLE+collection_id, qry, creds)
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
	resp, err := put(_DATA_PREAMBLE+collection_id, body, creds)
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
	query_map := query.serialize()
	query_bytes, err := json.Marshal(query_map)
	if err != nil {
		return err
	}
	qry := map[string]string{
		"query": url.QueryEscape(string(query_bytes)),
	}
	resp, err := delete(_DATA_PREAMBLE+collection_id, qry, creds)
	if err != nil {
		return fmt.Errorf("Error deleting data: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting data: %v", resp.Body)
	}
	return nil
}
