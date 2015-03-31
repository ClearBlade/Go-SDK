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
	if u.UserToken != "" {
		return [][]string{
			[]string{
				_USER_HEADER_KEY,
				u.UserToken,
			},
		}, nil
	} else if u.SystemSecret != "" && u.SystemKey != "" {
		return [][]string{
			[]string{
				_HEADER_SECRET_KEY,
				u.SystemSecret,
			},
			[]string{
				_HEADER_KEY_KEY,
				u.SystemKey,
			},
		}, nil
	} else {
		return [][]string{}, errors.New("No SystemSecret/SystemKey combo, or UserToken found")
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

func (d *DevClient) GetUserColumns(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(_USER_ADMIN+"/"+systemKey+"/columns", nil, creds)
	if err != nil {
		return nil, fmt.Errorf("Error getting user columns: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting user columns: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) CreateUserColumn(systemKey, columnName, columnType string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"column_name": columnName,
		"type":        columnType,
	}

	resp, err := post(_USER_ADMIN+"/"+systemKey+"/columns", data, creds)
	if err != nil {
		return fmt.Errorf("Error creating user column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error creating user column: %v", resp.Body)
	}

	return nil
}
