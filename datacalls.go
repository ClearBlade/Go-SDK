package GoSDK

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *Client) InsertData(collection_id string, data interface{}) error {
	resp, err := c.Post("/api/v/1/data/"+collection_id, data)
	if err != nil {
		return fmt.Errorf("Error inserting: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error inserting: %v", resp.Body)
	}
	return nil
}

func (c *Client) GetData(collection_id string, query [][]map[string]interface{}) (map[string]interface{}, error) {
	var qry map[string]string
	if query != nil {
		b, jsonErr := json.Marshal(query)
		if jsonErr != nil {
			return nil, fmt.Errorf("JSON Encoding error: %v", jsonErr)
		}
		qryStr := url.QueryEscape(string(b))
		qry = map[string]string{"query": qryStr}
	} else {
		qry = nil
	}
	resp, err := c.Get("/api/v/1/data/"+collection_id, qry)
	if err != nil {
		return nil, fmt.Errorf("Error getting data: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting data: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (c *Client) UpdateData(collection_id string, query [][]map[string]interface{}, changes map[string]interface{}) error {
	body := map[string]interface{}{
		"query": query,
		"$set":  changes,
	}
	resp, err := c.Put("/api/v/1/data/"+collection_id, body)
	if err != nil {
		return fmt.Errorf("Error updating data: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating data: %v", resp.Body)
	}
	return nil
}

func (c *Client) DeleteData(collection_id string, query [][]map[string]interface{}) error {
	var qry map[string]string
	if query != nil {
		b, jsonErr := json.Marshal(query)
		if jsonErr != nil {
			return fmt.Errorf("JSON Encoding error: %v", jsonErr)
		}
		qryStr := url.QueryEscape(string(b))
		qry = map[string]string{"query": qryStr}
	} else {
		return fmt.Errorf("Must supply a query to delete")
	}
	resp, err := c.Delete("/api/v/1/data/"+collection_id, qry)
	if err != nil {
		return fmt.Errorf("Error deleting data: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting data: %v", resp.Body)
	}
	return nil
}
