package GoSDK

import (
	"encoding/json"
	"fmt"
)

const (
	_DATA_PREAMBLE      = "/api/v/1/data/"
	_DATA_NAME_PREAMBLE = "/api/v/1/collection/"
	_DATA_V2_PREAMBLE   = "/api/v/2"
	_DATA_V3_PREAMBLE   = "/api/v/3"
	_DATA_V4_PREAMBLE   = "/api/v/4/"
)

// Inserts data into the platform. The interface is either a map[string]interface{} representing a row, or a []map[string]interface{} representing many rows.
func (u *UserClient) InsertData(collection_id string, data interface{}) error {
	_, err := insertdata(u, collection_id, data)
	return err
}

// Inserts data into the platform. The interface is either a map[string]interface{} representing a row, or a []map[string]interface{} representing many rows.
func (d *DeviceClient) InsertData(collection_id string, data interface{}) error {
	_, err := insertdata(d, collection_id, data)
	return err
}

// Inserts data into the platform. The interface is either a map[string]interface{} representing a row, or a []map[string]interface{} representing many rows.
func (d *DevClient) InsertData(collection_id string, data interface{}) error {
	_, err := insertdata(d, collection_id, data)
	return err
}

// CreateData is an alias for InsertData, but returns a response value, it should be a slice of strings representing the item ids (if not using an external datastore)
func (d *DevClient) CreateData(collection_id string, data interface{}) ([]interface{}, error) {
	return insertdata(d, collection_id, data)
}

// CreateData is an alias for InsertData, but returns a response value, it should be a slice of strings representing the item ids (if not using an external datastore)
func (u *UserClient) CreateData(collection_id string, data interface{}) ([]interface{}, error) {
	return insertdata(u, collection_id, data)
}

// CreateData is an alias for InsertData, but returns a response value, it should be a slice of strings representing the item ids (if not using an external datastore)
func (d *DeviceClient) CreateData(collection_id string, data interface{}) ([]interface{}, error) {
	return insertdata(d, collection_id, data)
}

func insertdata(c cbClient, collection_id string, data interface{}) ([]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(c, _DATA_PREAMBLE+collection_id, data, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error inserting: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error inserting: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

// GetData performs a query against a collection. The query object is discussed elsewhere. If the query object is nil, then it will return all of the data.
// The return value is a key-value of the types. Note that due to the transport mechanism being JSON, ints will be turned into float64s.
func (u *UserClient) GetData(collection_id string, query *Query) (map[string]interface{}, error) {
	return getdata(u, collection_id, query)
}

// GetData performs a query against a collection. The query object is discussed elsewhere. If the query object is nil, then it will return all of the data.
// The return value is a key-value of the types. Note that due to the transport mechanism being JSON, ints will be turned into float64s.
func (d *DeviceClient) GetData(collection_id string, query *Query) (map[string]interface{}, error) {
	return getdata(d, collection_id, query)
}

// GetDataByName performs a query against a collection, using the collection's name, rather than the ID. The query object is discussed elsewhere. If the query object is nil, then it will return all of the data.
// The return value is a key-value of the types. Note that due to the transport mechanism being JSON, ints will be turned into float64s.
func (u *UserClient) GetDataByName(collectionName string, query *Query) (map[string]interface{}, error) {
	return getDataByName(u, u.SystemKey, collectionName, query)
}

// GetDataByName performs a query against a collection, using the collection's name, rather than the ID. The query object is discussed elsewhere. If the query object is nil, then it will return all of the data.
// The return value is a key-value of the types. Note that due to the transport mechanism being JSON, ints will be turned into float64s.
func (d *DeviceClient) GetDataByName(collectionName string, query *Query) (map[string]interface{}, error) {
	return getDataByName(d, d.SystemKey, collectionName, query)
}

// GetDataByName performs a query against a collection, using the collection's name, rather than the ID. The query object is discussed elsewhere. If the query object is nil, then it will return all of the data.
// The return value is a key-value of the types. Note that due to the transport mechanism being JSON, ints will be turned into float64s.
func (d *DevClient) GetDataByName(collectionName string, query *Query) (map[string]interface{}, error) {
	return nil, fmt.Errorf("Developer cannot call GetDataByName, must use GetDataByNameWithSystemKey")
}

// GetData performs a query against a collection. The query object is discussed elsewhere. If the query object is nil, then it will return all of the data.
// The return value is a key-value of the types. Note that due to the transport mechanism being JSON, ints will be turned into float64s.
func (d *DevClient) GetData(collection_id string, query *Query) (map[string]interface{}, error) {
	return getdata(d, collection_id, query)
}

func (d *DevClient) GetDataTotal(collection_id string, query *Query) (map[string]interface{}, error) {
	return getdatatotal(d, collection_id, query)
}

func (d *DeviceClient) GetDataTotal(collection_id string, query *Query) (map[string]interface{}, error) {
	return getdatatotal(d, collection_id, query)
}

func (u *UserClient) GetDataTotal(collection_id string, query *Query) (map[string]interface{}, error) {
	return getdatatotal(u, collection_id, query)
}

func (d *DevClient) GetDataTotalByName(system_key, collection_name string, query *Query) (map[string]interface{}, error) {
	return getdatatotalbyname(d, system_key, collection_name, query)
}

func (d *DeviceClient) GetDataTotalByName(system_key, collection_name string, query *Query) (map[string]interface{}, error) {
	return getdatatotalbyname(d, system_key, collection_name, query)
}

func (u *UserClient) GetDataTotalByName(system_key, collection_name string, query *Query) (map[string]interface{}, error) {
	return getdatatotalbyname(u, system_key, collection_name, query)
}

func (u *UserClient) GetItemCount(collection_id string) (int, error) {
	return getItemCount(u, collection_id)
}

func (d *DeviceClient) GetItemCount(collection_id string) (int, error) {
	return getItemCount(d, collection_id)
}

func (d *DevClient) GetItemCount(collection_id string) (int, error) {
	return getItemCount(d, collection_id)
}

func getItemCount(c cbClient, collection_id string) (int, error) {
	creds, err := c.credentials()
	if err != nil {
		return -1, err
	}
	resp, err := get(c, _DATA_PREAMBLE+collection_id+"/count", nil, creds, nil)
	if err != nil {
		return -1, fmt.Errorf("Error getting count: %v", err)
	}
	if resp.StatusCode != 200 {
		return -1, fmt.Errorf("Error getting count: %v", resp.Body)
	}
	bod := resp.Body.(map[string]interface{})
	theCount := int(bod["count"].(float64))
	return theCount, nil

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
			"query": string(query_bytes),
		}
	} else {
		qry = nil
	}
	resp, err := get(c, _DATA_NAME_PREAMBLE+sysKey+"/"+collectionName, qry, creds, nil)
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
			"query": string(query_bytes),
		}
	} else {
		qry = nil
	}
	resp, err := get(c, _DATA_PREAMBLE+collection_id, qry, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting data: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting data: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func getdatatotal(c cbClient, collection_id string, query *Query) (map[string]interface{}, error) {

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
			"query": string(query_bytes),
		}
	} else {
		qry = nil
	}
	resp, err := get(c, _DATA_PREAMBLE+collection_id+"/count", qry, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting data: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting data: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func getdatatotalbyname(c cbClient, system_key, collection_name string, query *Query) (map[string]interface{}, error) {

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
			"query": string(query_bytes),
		}
	} else {
		qry = nil
	}
	resp, err := get(c, _DATA_V2_PREAMBLE+"/collection/"+system_key+"/"+collection_name+"/count", qry, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting data: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting data: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

//UpsertData mutates the values in extant rows, selecting them via a query. If the query is nil, it updates all rows
//changes should be a map of the names of the columns, and the value you want them updated to

func (u *UserClient) UpsertData(collection_id string, changes map[string]interface{}, conflictColumn string) (map[string]interface{}, error) {
	return upsertdata(u, collection_id, changes, conflictColumn)
}

func (d *DeviceClient) UpsertData(collection_id string, changes map[string]interface{}, conflictColumn string) (map[string]interface{}, error) {
	return upsertdata(d, collection_id, changes, conflictColumn)
}

// UpsertData mutates the values in extant rows, selecting them via a query. If the query is nil, it updates all rows
// changes should be a map of the names of the columns, and the value you want them updated to
func (d *DevClient) UpsertData(collection_id string, changes map[string]interface{}, conflictColumn string) (map[string]interface{}, error) {
	return upsertdata(d, collection_id, changes, conflictColumn)
}

//UpsertDataByName mutates the values in extant rows, selecting them via a query. If the query is nil, it updates all rows
//changes should be a map of the names of the columns, and the value you want them updated to

func (u *UserClient) UpsertDataByName(system_key, collection_name string, changes map[string]interface{}, conflictColumn string) (map[string]interface{}, error) {
	return upsertdataByName(u, system_key, collection_name, changes, conflictColumn)
}

func (d *DeviceClient) UpsertDataByName(system_key, collection_name string, changes map[string]interface{}, conflictColumn string) (map[string]interface{}, error) {
	return upsertdataByName(d, system_key, collection_name, changes, conflictColumn)
}

// UpsertDataByName mutates the values in extant rows, selecting them via a query. If the query is nil, it updates all rows
// changes should be a map of the names of the columns, and the value you want them updated to
func (d *DevClient) UpsertDataByName(system_key, collection_name string, changes map[string]interface{}, conflictColumn string) (map[string]interface{}, error) {
	return upsertdataByName(d, system_key, collection_name, changes, conflictColumn)
}

func upsertdata(c cbClient, collection_id string, data interface{}, conflictColumn string) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%sdata/%s/upsert?conflictColumn=%s", _DATA_V4_PREAMBLE, collection_id, conflictColumn)
	resp, err := put(c, url, data, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error inserting: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error inserting: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func upsertdataByName(c cbClient, systemKey, collectionName string, data interface{}, conflictColumn string) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%scollection/%s/%s/upsert?conflictColumn=%s", _DATA_V4_PREAMBLE, systemKey, collectionName, conflictColumn)
	resp, err := put(c, url, data, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error inserting: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error inserting: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

//UpdateData mutates the values in extant rows, selecting them via a query. If the query is nil, it updates all rows
//changes should be a map of the names of the columns, and the value you want them updated to

func (u *UserClient) UpdateData(collection_id string, query *Query, changes map[string]interface{}) error {
	return updatedata(u, collection_id, query, changes)
}

func (d *DeviceClient) UpdateData(collection_id string, query *Query, changes map[string]interface{}) error {
	return updatedata(d, collection_id, query, changes)
}

// UpdateData mutates the values in extant rows, selecting them via a query. If the query is nil, it updates all rows
// changes should be a map of the names of the columns, and the value you want them updated to
func (d *DevClient) UpdateData(collection_id string, query *Query, changes map[string]interface{}) error {
	return updatedata(d, collection_id, query, changes)
}

//UpdateDataByName mutates the values in extant rows, selecting them via a query. If the query is nil, it updates all rows
//changes should be a map of the names of the columns, and the value you want them updated to

func (u *UserClient) UpdateDataByName(system_key, collection_name string, query *Query, changes map[string]interface{}) (UpdateResponse, error) {
	return updatedataByName(u, system_key, collection_name, query, changes)
}

func (d *DeviceClient) UpdateDataByName(system_key, collection_name string, query *Query, changes map[string]interface{}) (UpdateResponse, error) {
	return updatedataByName(d, system_key, collection_name, query, changes)
}

// UpdateDataByName mutates the values in extant rows, selecting them via a query. If the query is nil, it updates all rows
// changes should be a map of the names of the columns, and the value you want them updated to
func (d *DevClient) UpdateDataByName(system_key, collection_name string, query *Query, changes map[string]interface{}) (UpdateResponse, error) {
	return updatedataByName(d, system_key, collection_name, query, changes)
}

func (u *UserClient) CreateDataByName(system_key, collection_name string, item interface{}) ([]interface{}, error) {
	return createDataByName(u, system_key, collection_name, item)
}

func (d *DeviceClient) CreateDataByName(system_key, collection_name string, item interface{}) ([]interface{}, error) {
	return createDataByName(d, system_key, collection_name, item)
}

func (d *DevClient) CreateDataByName(system_key, collection_name string, item interface{}) ([]interface{}, error) {
	return createDataByName(d, system_key, collection_name, item)
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
	resp, err := put(c, _DATA_PREAMBLE+collection_id, body, creds, nil)
	if err != nil {
		return fmt.Errorf("Error updating data: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating data: %v", resp.Body)
	}
	return nil
}

type UpdateResponse struct {
	Count float64
}

func updatedataByName(c cbClient, system_key, collection_name string, query *Query, changes map[string]interface{}) (UpdateResponse, error) {
	qry := query.serialize()
	body := map[string]interface{}{
		"query": qry,
		"$set":  changes,
	}
	creds, err := c.credentials()
	if err != nil {
		return UpdateResponse{}, err
	}
	resp, err := put(c, _DATA_NAME_PREAMBLE+system_key+"/"+collection_name, body, creds, nil)
	if err != nil {
		return UpdateResponse{}, fmt.Errorf("Error updating data: %v", err)
	}
	if resp.StatusCode != 200 {
		return UpdateResponse{}, fmt.Errorf("Error updating data: %v", resp.Body)
	}
	fmtBody := make(map[string]interface{})
	ok := true
	if fmtBody, ok = resp.Body.(map[string]interface{}); !ok {
		return UpdateResponse{}, fmt.Errorf("Unexpected response type from update. Body is - %+v\n", resp.Body)
	}
	if count, ok := fmtBody["count"].(float64); !ok {
		return UpdateResponse{}, fmt.Errorf("No count key in response type from update. Body is - %+v\n", fmtBody)
	} else {
		return UpdateResponse{
			Count: count,
		}, nil
	}
}

func createDataByName(c cbClient, system_key, collection_name string, item interface{}) ([]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(c, _DATA_NAME_PREAMBLE+system_key+"/"+collection_name, item, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error updating data: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error updating data: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

// DeleteData removes data from a collection according to what matches the query. If the query is nil, then all data will be removed.
func (u *UserClient) DeleteData(collection_id string, query *Query) error {
	return deletedata(u, collection_id, query)
}

// DeleteData removes data from a collection according to what matches the query. If the query is nil, then all data will be removed.
func (d *DeviceClient) DeleteData(collection_id string, query *Query) error {
	return deletedata(d, collection_id, query)
}

// DeleteData removes data from a collection according to what matches the query. If the query is nil, then all data will be removed.
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
			"query": string(query_bytes),
		}
	} else {
		qry = nil
	}
	resp, err := delete(c, _DATA_PREAMBLE+collection_id, qry, creds, nil)
	if err != nil {
		return fmt.Errorf("Error deleting data: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting data: %v", resp.Body)
	}
	return nil
}

// GetColumns gets a slice of map[string]interface{} of the column names and values.
// As map[string]interface{}{"ColumnName":"name","ColumnType":"typename in string", "PK":bool}
func (d *DevClient) GetColumns(collectionId, systemKey, systemSecret string) ([]interface{}, error) {
	return getColumns(d, collectionId, systemKey, systemSecret)
}

// GetColumns gets a slice of map[string]interface{} of the column names and values.
// As map[string]interface{}{"ColumnName":"name","ColumnType":"typename in string", "PK":bool}
func (u *UserClient) GetColumns(collection_id, systemKey, systemSecret string) ([]interface{}, error) {
	return getColumns(u, collection_id, "", "")
}

// GetColumns gets a slice of map[string]interface{} of the column names and values.
// As map[string]interface{}{"ColumnName":"name","ColumnType":"typename in string", "PK":bool}
func (d *DeviceClient) GetColumns(collection_id, systemKey, systemSecret string) ([]interface{}, error) {
	return getColumns(d, collection_id, "", "")
}

// GetColumnsByCollectionName gets a slice of map[string]interface{} of the column names and values.
// As map[string]interface{}{"ColumnName":"name","ColumnType":"typename in string", "PK":bool}
func (d *DevClient) GetColumnsByCollectionName(systemKey, collectionName string) ([]interface{}, error) {
	return getColumnsByCollectionName(d, systemKey, collectionName)
}

func getColumnsByCollectionName(c cbClient, systemKey, collectionName string) ([]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(c, _DATA_V2_PREAMBLE+"/collection/"+systemKey+"/"+collectionName+"/columns", nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting collection columns: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting collection columns: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

func getColumns(c cbClient, collection_id, systemKey, systemSecret string) ([]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}

	var headers map[string][]string = nil
	if systemKey != "" {
		headers = map[string][]string{
			"Clearblade-Systemkey":    []string{systemKey},
			"Clearblade-Systemsecret": []string{systemSecret},
		}
	}

	resp, err := get(c, _DATA_PREAMBLE+collection_id+"/columns", nil, creds, headers)
	if err != nil {
		return nil, fmt.Errorf("Error getting collection columns: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting collection columns: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

// GetAllCollections retrieves a list of every collection in the system
// The return value is a slice of strings
func (d *DevClient) GetAllCollections(systemKey string) ([]interface{}, error) {
	return getAllCollections(d, d.preamble()+"/allcollections", systemKey)
}

func (u *UserClient) GetAllCollections(systemKey string) ([]interface{}, error) {
	return getAllCollections(u, _DATA_V3_PREAMBLE+"/allcollections/"+systemKey, systemKey)
}

func getAllCollections(c cbClient, endpoint, systemKey string) ([]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(c, endpoint, map[string]string{
		"appid": systemKey,
	}, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error fetchings all collections: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error fetchings all collections %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) NewCollection(systemKey, name string) (string, error) {
	return createNewCollection(d, d.preamble(), name, systemKey)
}

func (u *UserClient) NewCollection(systemKey, name string) (string, error) {
	return createNewCollection(u, _DATA_V3_PREAMBLE, name, systemKey)
}

// CreateCollection creates a new collection
func createNewCollection(c cbClient, preamble, name, systemKey string) (string, error) {
	creds, err := c.credentials()
	if err != nil {
		return "", err
	}
	resp, err := post(c, preamble+"/collectionmanagement", map[string]interface{}{
		"name":  name,
		"appID": systemKey,
	}, creds, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating collection: %v", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error creating collection %v", resp.Body)
	}
	return resp.Body.(map[string]interface{})["collectionID"].(string), nil
}

// GetCollectionInfo retrieves some describing information on the specified collection
// Keys "name","collectoinID","appID", and much, much more!
func (d *DevClient) GetCollectionInfo(collection_id string) (map[string]interface{}, error) {
	return getCollectionInfo(d, d.preamble(), collection_id)
}

func (u *UserClient) GetCollectionInfo(collection_id string) (map[string]interface{}, error) {
	return getCollectionInfo(u, _DATA_V3_PREAMBLE, collection_id)
}

func getCollectionInfo(c cbClient, preamble, collection_id string) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return map[string]interface{}{}, err
	}
	resp, err := get(c, preamble+"/collectionmanagement", map[string]string{
		"id": collection_id,
	}, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting collection info: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting collection info: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

// AddColumn adds a column to a collection. Note that this does not apply to collections backed by a non-default datastore.
func (d *DevClient) AddColumn(collection_id, column_name, column_type string) error {
	return addColumn(d, d.preamble(), collection_id, column_name, column_type)
}

func (u *UserClient) AddColumn(collection_id, column_name, column_type string) error {
	return addColumn(u, _DATA_V3_PREAMBLE, collection_id, column_name, column_type)
}

func addColumn(c cbClient, preamble, collection_id, column_name, column_type string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	resp, err := put(c, preamble+"/collectionmanagement", map[string]interface{}{
		"id": collection_id,
		"addColumn": map[string]interface{}{
			"name": column_name,
			"type": column_type,
		},
	}, creds, nil)
	if err != nil {
		return fmt.Errorf("Error adding column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error adding column: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) ConvertCollectionToHypertable(systemKey, collectionName string, options map[string]interface{}) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := post(d, _DATA_V4_PREAMBLE+"/collection/"+systemKey+"/"+collectionName+"/hypertable", options, creds, nil)
	if err != nil {
		return fmt.Errorf("Error converting collection to hypertable: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error converting collection to hypertable: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) UpdateHypertableProperties(systemKey, collectionName string, options map[string]interface{}) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put(d, _DATA_V4_PREAMBLE+"/collection/"+systemKey+"/"+collectionName+"/hypertable", options, creds, nil)
	if err != nil {
		return fmt.Errorf("Error updating hypertable properties: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating hypertable properties: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) DropHypertableChunks(systemKey, collectionName string, body map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := deleteWithBody(d, _DATA_V4_PREAMBLE+"/collection/"+systemKey+"/"+collectionName+"/hypertable", body, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error deleting hypertable chunks: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error deleting hypertable chunks: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetAllContinuousAggregatesForCollection(systemKey, collectionName string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _DATA_V4_PREAMBLE+"/collection/"+systemKey+"/"+collectionName+"/listcontinuousaggregates", nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting all continuous aggregates for collection: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting all continuous aggregates for collection: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetContinuousAggregate(systemKey, collectionName, aggregateName string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _DATA_V4_PREAMBLE+"/collection/"+systemKey+"/"+collectionName+"/"+aggregateName+"/continuousaggregate", nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting continuous aggregate: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting continuous aggregate: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) CreateContinuousAggregate(systemKey, collectionName, aggregateName string, properties map[string]interface{}) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := post(d, _DATA_V4_PREAMBLE+"/collection/"+systemKey+"/"+collectionName+"/"+aggregateName+"/continuousaggregate", properties, creds, nil)
	if err != nil {
		return fmt.Errorf("Error creating continuous aggregate: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error creating continuous aggregate: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) UpdateContinuousAggregate(systemKey, collectionName, aggregateName string, properties map[string]interface{}) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put(d, _DATA_V4_PREAMBLE+"/collection/"+systemKey+"/"+collectionName+"/"+aggregateName+"/continuousaggregate", properties, creds, nil)
	if err != nil {
		return fmt.Errorf("Error updating continuous aggregate: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating continuous aggregate: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) DeleteContinuousAggregate(systemKey, collectionName, aggregateName string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, _DATA_V4_PREAMBLE+"/collection/"+systemKey+"/"+collectionName+"/"+aggregateName+"/continuousaggregate", nil, creds, nil)
	if err != nil {
		return fmt.Errorf("Error deleting continuous aggregate: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting continuous aggregate: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) QueryContinuousAggregate(systemKey, collectionName, aggregateName string, query *Query) (map[string]interface{}, error) {
	creds, err := d.credentials()
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
			"query": string(query_bytes),
		}
	} else {
		qry = nil
	}
	resp, err := get(d, _DATA_V4_PREAMBLE+"/collection/"+systemKey+"/"+collectionName+"/"+aggregateName+"/continuousaggregate/query", qry, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error querying continuous aggregate: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error querying continuous aggregate: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

// DeleteColumn removes a column from a collection. Note that this does not apply to collections backed by a non-default datastore.
func (d *DevClient) DeleteColumn(collection_id, column_name string) error {
	return deleteColumn(d, d.preamble(), collection_id, column_name)
}

func (u *UserClient) DeleteColumn(collection_id, column_name string) error {
	return deleteColumn(u, _DATA_V3_PREAMBLE, collection_id, column_name)
}

func deleteColumn(c cbClient, preamble, collection_id, column_name string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	resp, err := put(c, preamble+"/collectionmanagement", map[string]interface{}{
		"id":           collection_id,
		"deleteColumn": column_name,
	}, creds, nil)
	if err != nil {
		return fmt.Errorf("Error deleting column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting column: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) RenameColumn(collectionID, oldColumnName, newColumnName string) error {
	return renameColumn(d, d.preamble(), collectionID, oldColumnName, newColumnName)
}

func (u *UserClient) RenameColumn(collectionID, oldColumnName, newColumnName string) error {
	return renameColumn(u, _DATA_V3_PREAMBLE, collectionID, oldColumnName, newColumnName)
}

func renameColumn(c cbClient, preamble, collectionID, oldColumnName, newColumnName string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	resp, err := put(c, preamble+"/collectionmanagement", map[string]interface{}{
		"id": collectionID,
		"renameColumn": map[string]string{
			"from": oldColumnName,
			"to":   newColumnName,
		},
	}, creds, nil)
	if err != nil {
		return fmt.Errorf("failed to rename column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to rename column: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) RenameCollection(collectionID, newName string) error {
	return renameCollection(d, d.preamble(), collectionID, newName)
}

func (u *UserClient) RenameCollection(collectionID, newName string) error {
	return renameCollection(u, _DATA_V3_PREAMBLE, collectionID, newName)
}

func renameCollection(c cbClient, preamble, collectionID, newName string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	resp, err := put(c, preamble+"/collectionmanagement", map[string]interface{}{
		"id":               collectionID,
		"renameCollection": newName,
	}, creds, nil)
	if err != nil {
		return fmt.Errorf("failed to rename collection: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to rename collection: %v", resp.Body)
	}
	return nil
}

// DeleteCollection deletes the collection. Note that this does not apply to collections backed by a non-default datastore.
func (d *DevClient) DeleteCollection(colID string) error {
	return deleteCollection(d, d.preamble(), colID)
}

func (u *UserClient) DeleteCollection(colID string) error {
	return deleteCollection(u, _DATA_V3_PREAMBLE, colID)
}

func deleteCollection(c cbClient, preamble, colID string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(c, preamble+"/collectionmanagement", map[string]string{
		"id": colID,
	}, creds, nil)
	if err != nil {
		return fmt.Errorf("Error deleting collection %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting collection %v", resp.Body)
	}
	return nil
}

func (d *DevClient) GetDataByNameWithSystemKey(systemKey, collectionName string, query *Query) (map[string]interface{}, error) {
	return getDataByName(d, systemKey, collectionName, query)
}

func (u *UserClient) GetDataByNameWithSystemKey(systemKey, collectionName string, query *Query) (map[string]interface{}, error) {
	return getDataByName(u, systemKey, collectionName, query)
}

func (d *DeviceClient) GetDataByNameWithSystemKey(systemKey, collectionName string, query *Query) (map[string]interface{}, error) {
	return getDataByName(d, systemKey, collectionName, query)
}

func (u *UserClient) CreateIndex(systemKey, collectionName, columnToIndex string) error {
	return createIndex(u, systemKey, collectionName, columnToIndex)
}

func (d *DeviceClient) CreateIndex(systemKey, collectionName, columnToIndex string) error {
	return createIndex(d, systemKey, collectionName, columnToIndex)
}

func (d *DevClient) CreateIndex(systemKey, collectionName, columnToIndex string) error {
	return createIndex(d, systemKey, collectionName, columnToIndex)
}

func createIndex(c cbClient, systemKey, collectionName, columnToIndex string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%scollection/%s/%s/index?columnName=%s", _DATA_V4_PREAMBLE, systemKey, collectionName, columnToIndex)
	resp, err := post(c, url, nil, creds, nil)
	if err != nil {
		return fmt.Errorf("Error sending request for creating index: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error creating index: %v", resp.Body)
	}
	return nil
}

func (u *UserClient) CreateUniqueIndex(systemKey, collectionName, columnToIndex string) error {
	return createUniqueIndex(u, systemKey, collectionName, columnToIndex)
}

func (d *DeviceClient) CreateUniqueIndex(systemKey, collectionName, columnToIndex string) error {
	return createUniqueIndex(d, systemKey, collectionName, columnToIndex)
}

func (d *DevClient) CreateUniqueIndex(systemKey, collectionName, columnToIndex string) error {
	return createUniqueIndex(d, systemKey, collectionName, columnToIndex)
}

func createUniqueIndex(c cbClient, systemKey, collectionName, columnToIndex string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%scollection/%s/%s/uniqueindex?columnName=%s", _DATA_V4_PREAMBLE, systemKey, collectionName, columnToIndex)
	resp, err := post(c, url, nil, creds, nil)
	if err != nil {
		return fmt.Errorf("Error sending request for creating unique index: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error creating unique index: %v", resp.Body)
	}
	return nil
}

func (u *UserClient) DropUniqueIndex(systemKey, collectionName, columnToIndex string) error {
	return dropUniqueIndex(u, systemKey, collectionName, columnToIndex)
}

func (d *DeviceClient) DropUniqueIndex(systemKey, collectionName, columnToIndex string) error {
	return dropUniqueIndex(d, systemKey, collectionName, columnToIndex)
}

func (d *DevClient) DropUniqueIndex(systemKey, collectionName, columnToIndex string) error {
	return dropUniqueIndex(d, systemKey, collectionName, columnToIndex)
}

func dropUniqueIndex(c cbClient, systemKey, collectionName, columnToIndex string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%scollection/%s/%s/uniqueindex?columnName=%s", _DATA_V4_PREAMBLE, systemKey, collectionName, columnToIndex)
	resp, err := delete(c, url, nil, creds, nil)
	if err != nil {
		return fmt.Errorf("Error sending request for dropping unique index: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error dropping unique index: %v", resp.Body)
	}
	return nil
}

func (u *UserClient) DropIndex(systemKey, collectionName, columnToIndex string) error {
	return dropIndex(u, systemKey, collectionName, columnToIndex)
}

func (d *DeviceClient) DropIndex(systemKey, collectionName, columnToIndex string) error {
	return dropIndex(d, systemKey, collectionName, columnToIndex)
}

func (d *DevClient) DropIndex(systemKey, collectionName, columnToIndex string) error {
	return dropIndex(d, systemKey, collectionName, columnToIndex)
}

func dropIndex(c cbClient, systemKey, collectionName, columnToIndex string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%scollection/%s/%s/index?columnName=%s", _DATA_V4_PREAMBLE, systemKey, collectionName, columnToIndex)
	resp, err := delete(c, url, nil, creds, nil)
	if err != nil {
		return fmt.Errorf("Error sending request for dropping index: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error dropping index: %v", resp.Body)
	}
	return nil
}

func (u *UserClient) ListIndexes(systemKey, collectionName string) (map[string]interface{}, error) {
	return listIndexes(u, systemKey, collectionName)
}

func (d *DeviceClient) ListIndexes(systemKey, collectionName string) (map[string]interface{}, error) {
	return listIndexes(d, systemKey, collectionName)
}

func (d *DevClient) ListIndexes(systemKey, collectionName string) (map[string]interface{}, error) {
	return listIndexes(d, systemKey, collectionName)
}

func listIndexes(c cbClient, systemKey, collectionName string) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%scollection/%s/%s/listindexes", _DATA_V4_PREAMBLE, systemKey, collectionName)
	resp, err := get(c, url, nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error sending request for list indexes: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error listing indexes: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (u *UserClient) ListIndexesWithID(collectionID string) (map[string]interface{}, error) {
	return listIndexeswithid(u, collectionID)
}

func (d *DeviceClient) ListIndexesWithID(collectionID string) (map[string]interface{}, error) {
	return listIndexeswithid(d, collectionID)
}

func (d *DevClient) ListIndexesWithID(collectionID string) (map[string]interface{}, error) {
	return listIndexeswithid(d, collectionID)
}

func listIndexeswithid(c cbClient, collectionID string) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%sdata/%s/listindexes", _DATA_V4_PREAMBLE, collectionID)
	resp, err := get(c, url, nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error sending request for list indexes: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error listing indexes: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) RawExec(systemKey, query string, params []interface{}) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	url := _DATA_V4_PREAMBLE + "database/" + systemKey + "/exec"
	data := map[string]interface{}{"query": query, "params": params}
	_, err = post(d, url, data, creds, nil)
	if err != nil {
		return fmt.Errorf("Error executing %v with args %v: %v", query, params, err)
	}
	return nil
}

func (d *DevClient) RawQuery(systemKey, query string, params []interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	url := _DATA_V4_PREAMBLE + "database/" + systemKey + "/query"
	data := map[string]interface{}{"query": query, "params": params}
	resp, err := post(d, url, data, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error executing %v with args %v: %v", query, params, err)
	}
	return resp.Body.(map[string]interface{}), nil
}
