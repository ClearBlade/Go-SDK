package GoSDK

import (
	"errors"
	"fmt"
)

const (
	_USER_HEADER_KEY = "ClearBlade-UserToken"
	_USER_PREAMBLE   = "/api/v/1/user"
	_USER_ADMIN      = "/admin/user"
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

func (u *UserClient) getMessageId() uint16 {
	return uint16(u.mrand.Int())
}

//GetUserColumns returns the description of the columns in the user table
//Returns a structure shaped []map[string]interface{}{map[string]interface{}{"ColumnName":"blah","ColumnType":"int"}}
func (d *DevClient) GetUserColumns(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(_USER_ADMIN+"/"+systemKey+"/columns", nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting user columns: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting user columns: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

//CreateUserColumn creates a new column in the user table
func (d *DevClient) CreateUserColumn(systemKey, columnName, columnType string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"column_name": columnName,
		"type":        columnType,
	}

	resp, err := post(_USER_ADMIN+"/"+systemKey+"/columns", data, creds, nil)
	if err != nil {
		return fmt.Errorf("Error creating user column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error creating user column: %v", resp.Body)
	}

	return nil
}
