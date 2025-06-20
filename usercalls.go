package GoSDK

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	_USER_HEADER_KEY = "ClearBlade-UserToken"
	_USER_PREAMBLE   = "/api/v/1/user"
	_USER_V2         = "/api/v/2/user"
	_USER_V4         = "/api/v/4/user"
	_USER_ADMIN      = "/admin/user"
	_USER_SESSION    = "/admin/v/4/session"
)

func (u *UserClient) credentials() ([][]string, error) {
	ret := make([][]string, 0)
	if u.UserToken != "" {
		ret = append(ret, []string{
			_USER_HEADER_KEY,
			u.UserToken,
		})
	}
	if u.SystemSecret != "" && u.SystemKey != "" {
		ret = append(ret, []string{
			_HEADER_SECRET_KEY,
			u.SystemSecret,
		})
		ret = append(ret, []string{
			_HEADER_KEY_KEY,
			u.SystemKey,
		})

	}

	if len(ret) == 0 {
		return [][]string{}, errors.New("No SystemSecret/SystemKey combo, or UserToken found")
	} else {
		return ret, nil
	}
}

func (u *UserClient) preamble() string {
	return _USER_PREAMBLE
}

func (u *UserClient) getSystemInfo() (string, string) {
	return u.SystemKey, u.SystemSecret
}

func (u *UserClient) setToken(t string) {
	u.UserToken = t
}
func (u *UserClient) getToken() string {
	return u.UserToken
}
func (u *UserClient) getRefreshToken() string {
	return u.RefreshToken
}
func (u *UserClient) setRefreshToken(body map[string]interface{}) {
	u.RefreshToken = nicelySetRefreshToken(body)
}
func (u *UserClient) setExpiresAt(t float64) {
	u.ExpiresAt = t
}
func (u *UserClient) getExpiresAt() float64 {
	return u.ExpiresAt
}

func (u *UserClient) getMessageId() uint16 {
	return uint16(time.Now().UnixNano())
}

func (u *UserClient) GetUserCount(systemKey string) (int, error) {
	creds, err := u.credentials()
	if err != nil {
		return -1, err
	}
	resp, err := get(u, u.preamble()+"/count", nil, creds, nil)
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

func (d *DevClient) RegisterNewUserWithOptions(systemKey, systemSecret string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("/admin/user/%s", systemKey)
	headers := make(map[string][]string)
	headers["Clearblade-Systemkey"] = []string{systemKey}
	headers["Clearblade-Systemsecret"] = []string{systemSecret}
	resp, err := post(d, endpoint, data, creds, headers)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error registering user: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) GetUserCountWithQuery(systemKey string, query *Query) (CountResp, error) {
	creds, err := d.credentials()
	if err != nil {
		return CountResp{Count: 0}, err
	}

	qry, err := createQueryMap(query)
	if err != nil {
		return CountResp{Count: 0}, err
	}

	resp, err := get(d, _USER_ADMIN+"/"+systemKey+"/count", qry, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return CountResp{Count: 0}, err
	}
	rval, ok := resp.Body.(map[string]interface{})
	if !ok {
		return CountResp{Count: 0}, fmt.Errorf("Bad type returned by getDevicesCount: %T, %s", resp.Body, resp.Body.(string))
	}

	return CountResp{
		Count: rval["count"].(float64),
	}, nil
}

func (d *DevClient) GetUsersWithQuery(systemKey string, query *Query) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	qry, err := createQueryMap(query)
	if err != nil {
		return nil, err
	}

	resp, err := get(d, _USER_ADMIN+"/"+systemKey, qry, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

// GetUserColumns returns the description of the columns in the user table
// Returns a structure shaped []map[string]interface{}{map[string]interface{}{"ColumnName":"blah","ColumnType":"int"}}
func (d *DevClient) GetUserColumns(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(d, _USER_ADMIN+"/"+systemKey+"/columns", nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting user columns: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting user columns: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

// CreateUserColumn creates a new column in the user table
func (d *DevClient) CreateUserColumn(systemKey, columnName, columnType string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"column_name": columnName,
		"type":        columnType,
	}

	resp, err := post(d, _USER_ADMIN+"/"+systemKey+"/columns", data, creds, nil)
	if err != nil {
		return fmt.Errorf("Error creating user column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error creating user column: %v", resp.Body)
	}

	return nil
}

func (d *DevClient) DeleteUserColumn(systemKey, columnName string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	data := map[string]string{"column": columnName}

	resp, err := delete(d, _USER_ADMIN+"/"+systemKey+"/columns", data, creds, nil)
	if err != nil {
		return fmt.Errorf("Error deleting user column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting user column: %v", resp.Body)
	}

	return nil
}

func (u *UserClient) UpdateUser(userQuery *Query, changes map[string]interface{}) error {
	return updateUser(u, userQuery, changes)
}

func updateUser(c cbClient, userQuery *Query, changes map[string]interface{}) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	query := userQuery.serialize()
	body := map[string]interface{}{
		"query":   query,
		"changes": changes,
	}

	resp, err := put(c, _USER_V2+"/info", body, creds, nil)
	if err != nil {
		return fmt.Errorf("Error updating data: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating data: %v", resp.Body)
	}

	return nil
}

func (d *DevClient) GetUserSession(systemKey string, query *Query) ([]interface{}, error) {
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
	resp, err := get(d, _USER_SESSION+"/"+systemKey+"/user", qry, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting user session data: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting user session data: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) DeleteUserSession(systemKey string, query *Query) error {
	creds, err := d.credentials()
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
	resp, err := delete(d, _USER_SESSION+"/"+systemKey+"/user", qry, creds, nil)
	if err != nil {
		return fmt.Errorf("Error deleting user session data: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting user session data: %v", resp.Body)
	}
	return nil
}

func (u *UserClient) UpdateUserPassword(userID, newPassword string) error {
	creds, err := u.credentials()
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"user": userID,
		"changes": map[string]interface{}{
			"password": newPassword,
		},
	}

	resp, err := put(u, _USER_V4+"/manage", body, creds, nil)
	if err != nil {
		return fmt.Errorf("Error updating password: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating password: %v", resp.Body)
	}

	return nil
}

type RoleChanges struct {
	Add    []string `json:"add"`
	Delete []string `json:"delete"`
}

func (u *UserClient) UpdateUserRoles(userID string, changes RoleChanges) error {
	creds, err := u.credentials()
	if err != nil {
		return err
	}
	body := map[string]interface{}{
		"user": userID,
		"changes": map[string]interface{}{
			"roles": changes,
		},
	}

	resp, err := put(u, _USER_V4+"/manage", body, creds, nil)
	if err != nil {
		return fmt.Errorf("Error updating roles: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating roles: %v", resp.Body)
	}

	return nil
}

func (u *UserClient) GetOwnUserInfo() (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(u, u.preamble()+"/info", nil, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting own user info: %v", resp.Body)
	}
	rawData, ok := resp.Body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Error parsing response")
	}
	return rawData, nil
}

func (u *UserClient) GetUserInfo(systemKey, email string) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	query := NewQuery()
	query.EqualTo("email", email)
	var qry map[string]string
	query_map := query.serialize()
	query_bytes, err := json.Marshal(query_map)
	if err != nil {
		return nil, err
	}
	qry = map[string]string{
		"query": string(query_bytes),
	}
	resp, err := get(u, u.preamble(), qry, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting user %s: %v", email, resp.Body)
	}
	rawData, ok := resp.Body.(map[string]interface{})["Data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Error parsing response")
	}
	if len(rawData) == 0 {
		return nil, fmt.Errorf("User with email %s does not exist", email)
	}
	if len(rawData) != 1 {
		return nil, fmt.Errorf("Got more than one user for email %s", email)
	}

	return rawData[0].(map[string]interface{}), nil
}

func (u *UserClient) GetAllUsers(systemKey string) ([]map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	resp, err := get(u, u.preamble(), nil, creds, nil)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting all users: %v", resp.Body)
	}
	dbResponse := resp.Body.(map[string]interface{})
	rawData := dbResponse["Data"].([]interface{})

	rval := make([]map[string]interface{}, len(rawData))
	for idx, oneRsp := range rawData {
		rval[idx] = oneRsp.(map[string]interface{})
	}

	return rval, nil
}

func (u *UserClient) GetUserRoles(systemKey, userId string) ([]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(u, fmt.Sprintf("/api/v/3/user/roles/%s?user=%s", systemKey, userId), nil, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting user roles: %v", resp.Body)
	}

	return resp.Body.([]interface{}), nil
}

func ConnectedUsers(client cbClient, systemKey string) (map[string]interface{}, error) {
	creds, err := client.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(client, _USER_PREAMBLE+"/"+systemKey+"/connections", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func UserConnections(client cbClient, systemKey, deviceName string) (map[string]interface{}, error) {
	creds, err := client.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(client, _USER_PREAMBLE+"/"+systemKey+"/connections/"+deviceName, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func ConnectedUserCount(client cbClient, systemKey string) (map[string]interface{}, error) {
	creds, err := client.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(client, _USER_PREAMBLE+"/"+systemKey+"/connectioncount", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}
