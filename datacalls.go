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

//Inserts data into the platform. The interface is either a map[string]interface{} representing a row, or a []map[string]interface{} representing many rows.
func (u *UserClient) InsertData(collection_id string, data interface{}) error {
	_, err := insertdata(u, collection_id, data)
	return err
}

//Inserts data into the platform. The interface is either a map[string]interface{} representing a row, or a []map[string]interface{} representing many rows.
func (d *DevClient) InsertData(collection_id string, data interface{}) error {
	_, err := insertdata(d, collection_id, data)
	return err
}

//CreateData is an alias for InsertData, but returns a response value
func (d *DevClient) CreateData(collection_id string, data interface{}) ([]interface{}, error) {
	resp, err := insertdata(d, collection_id, data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//CreateData is an alias for InsertData, but returns a response value
func (u *UserClient) CreateData(collection_id string, data interface{}) ([]interface{}, error) {
	resp, err := insertdata(u, collection_id, data)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func insertdata(c cbClient, collection_id string, data interface{}) ([]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(_DATA_PREAMBLE+collection_id, data, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error inserting: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error inserting: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

//GetData performs a query against a collection. The query object is discussed elsewhere. If the query object is nil, then it will return all of the data.
func (u *UserClient) GetData(collection_id string, query *Query) (map[string]interface{}, error) {
	return getdata(u, collection_id, query)
}

//GetDataByName performs a query against a collection, using the collection's name, rather than the ID. The query object is discussed elsewhere. If the query object is nil, then it will return all of the data.
func (u *UserClient) GetDataByName(collectionName string, query *Query) (map[string]interface{}, error) {
	return getDataByName(u, u.SystemKey, collectionName, query)
}

//GetDataByName performs a query against a collection, using the collection's name, rather than the ID. The query object is discussed elsewhere. If the query object is nil, then it will return all of the data.
func (d *DevClient) GetDataByName(collectionName string, query *Query) (map[string]interface{}, error) {
	return nil, fmt.Errorf("Developer cannot call this yet")
}

//GetData performs a query against a collection. The query object is discussed elsewhere. If the query object is nil, then it will return all of the data.
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

//UpdateData mutates the values in extant rows, selecting them via a query. If the query is nil, it updates all rows
//changes should be a map of the names of the columns, and the value you want them updated to
func (u *UserClient) UpdateData(collection_id string, query *Query, changes map[string]interface{}) error {
	err := updatedata(u, collection_id, query, changes)
	return err
}

//UpdateData mutates the values in extant rows, selecting them via a query. If the query is nil, it updates all rows
//changes should be a map of the names of the columns, and the value you want them updated to
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

//DeleteData removes data from a collection according to what matches the query. If the query is nil, then all data will be removed.
func (u *UserClient) DeleteData(collection_id string, query *Query) error {
	return deletedata(u, collection_id, query)
}

//DeleteData removes data from a collection according to what matches the query. If the query is nil, then all data will be removed.
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

//GetColumns gets a slice of map[string]interface{} of the column names and values.
//As map[string]interface{}{"ColumnName":"name","ColumnType":"typename in string", "PK":bool}
func (d *DevClient) GetColumns(collection_id string) ([]interface{}, error) {
	return getColumns(d, collection_id)
}

//GetColumns gets a slice of map[string]interface{} of the column names and values.
//As map[string]interface{}{"ColumnName":"name","ColumnType":"typename in string", "PK":bool}
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

//GetDataByKeyAndName is unimplemented
func (d *DevClient) GetDataByKeyAndName(string, string, *Query) (map[string]interface{}, error) {
	return nil, fmt.Errorf("Unimplemented")
}

//GetDataByKeyAndName is unimplemented
func (d *UserClient) GetDataByKeyAndName(string, string, *Query) (map[string]interface{}, error) {
	return nil, fmt.Errorf("Unimplemented")

}
