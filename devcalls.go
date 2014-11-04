package GoSDK

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	_DEV_HEADER_KEY = "ClearBlade-DevToken"
	_DEV_PREAMBLE   = "admin"
)

type System struct {
	Key         string
	Secret      string
	Name        string
	Description string
	Users       bool
}

func (d *DevClient) NewSystem(name, description string, users bool) (string, error) {
	resp, err := post("/admin/systemmanagement", map[string]interface{}{
		"name":          name,
		"description":   description,
		"auth_required": users,
	}, d.creds())
	if err != nil {
		return "", fmt.Errorf("Error creating new system: %v", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error Creating new system: %v", resp.Body)
	}

	// TODO we need to make this json
	return strings.TrimSpace(strings.Split(resp.Body.(string), ":")[1]), nil
}

func (d *DevClient) GetSystem(key string) (*System, error) {
	//TODO:REPLACE
	// if _, ok := c.Headers["ClearBlade-DevToken"]; !ok {
	// 	return nil, fmt.Errorf("Error getting system: No DevToken Supplied")
	// }
	sysResp, sysErr := get("/admin/systemmanagement", map[string]string{"id": key}, d.creds())
	if sysErr != nil {
		return nil, fmt.Errorf("Error gathering system information: %v", sysErr)
	}
	if sysResp.StatusCode != 200 {
		return nil, fmt.Errorf("Error gathering system information: %v", sysResp.Body)
	}
	sysMap, isMap := sysResp.Body.(map[string]interface{})
	if !isMap {
		return nil, fmt.Errorf("Error gathering system information: incorrect return type\n")
	}
	newSys := &System{
		Key:         sysMap["appID"].(string),
		Secret:      sysMap["appSecret"].(string),
		Name:        sysMap["name"].(string),
		Description: sysMap["description"].(string),
		Users:       sysMap["auth_required"].(bool),
	}
	return newSys, nil

}

func (d *DevClient) DeleteSystem(s string) error {
	resp, err := delete("/admin/systemmanagement", map[string]string{"id": s}, d.creds())
	if err != nil {
		return fmt.Errorf("Error deleting system: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting system: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) SetSystemName(system_key, system_name string) error {
	resp, err := put("/admin/systemmanagement", map[string]interface{}{
		"id":   system_key,
		"name": system_name,
	}, d.creds())
	if err != nil {
		return fmt.Errorf("Error changing system name: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error changing system name: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) SetSystemDescription(system_key, system_description string) error {
	resp, err := c.Put("/admin/systemmanagement", map[string]interface{}{
		"id":          system_key,
		"description": system_description,
	}, d.creds())
	if err != nil {
		return fmt.Errorf("Error changing system description: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error changing system description: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) SetSystemAuthOn(system_key string) error {
	resp, err := put("/admin/systemmanagement", map[string]interface{}{
		"id":            system_key,
		"auth_required": true,
	}, d.creds())
	if err != nil {
		return fmt.Errorf("Error changing system auth: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error changing system auth: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) SetSystemAuthOff(system_key string) error {
	resp, err := put("/admin/systemmanagement", map[string]interface{}{
		"id":            system_key,
		"auth_required": false,
	}, d.creds())
	if err != nil {
		return fmt.Errorf("Error changing system auth: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error changing system auth: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) DevUserInfo() error {
	resp, err := get("/admin/userinfo", nil, c.creds())
	if err != nil {
		return fmt.Errorf("Error getting userdata: %v", err)
	}
	log.Printf("HERE IS THE BODY: %+v\n", resp)
	return nil
}

func (d *DevClient) NewCollection(systemKey, name string) (string, error) {
	resp, err := post("/admin/collectionmanagement", map[string]interface{}{
		"name":  name,
		"appID": systemKey,
	}, d.creds())
	if err != nil {
		return "", fmt.Errorf("Error creating collection: %v", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error creating collection %v", resp.Body)
	}
	return resp.Body.(map[string]interface{})["collectionID"].(string), nil
}

func (d *DevClient) DeleteCollection(colId string) error {
	resp, err := c.Delete("/admin/collectionmanagement", map[string]string{
		"id": colId,
	}, d.creds())
	if err != nil {
		return fmt.Errorf("Error deleting collection %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting collection %v", resp.Body)
	}
	return nil
}

func (d *DevClient) AddColumn(collection_id, column_name, column_type string) error {
	resp, err := put("/admin/collectionmanagement", map[string]interface{}{
		"id": collection_id,
		"addColumn": map[string]interface{}{
			"name": column_name,
			"type": column_type,
		},
	}, d.creds())
	if err != nil {
		return fmt.Errorf("Error adding column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error adding column: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) DeleteColumn(collection_id, column_name string) error {
	resp, err := put("/admin/collectionmanagement", map[string]interface{}{
		"id":           collection_id,
		"deleteColumn": column_name,
	}, d.creds())
	if err != nil {
		return fmt.Errorf("Error deleting column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting column: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) GetCollectionInfo(collection_id string) (map[string]interface{}, error) {
	resp, err := c.Get("/admin/collectionmanagement", map[string]string{
		"id": collection_id,
	}, d.creds())
	if err != nil {
		return nil, fmt.Errorf("Error getting collection info: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting collection info: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

//second verse, same as the first, eh?
func (d *DevClient) creds() ([][]string, error) {
	if d.DevToken != "" {
		return [][]string{
			[]string{
				_DEV_HEADER_KEY,
				d.DevToken,
			},
		}, nil
	} else if d.SystemSecret != "" && d.SystemKey != "" {
		return [][]string{
			[]string{
				_HEADER_SECRET_KEY,
				d.SystemSecret,
			},
			[]string{
				_HEADER_KEY_KEY,
				d.SystemKey,
			},
		}, nil
	} else {
		return [][]string{}, errors.New("No SystemSecret/SystemKey combo, or UserToken found")
	}
}

func (d *DevClient) preamble() string {
	return _DEV_PREAMBLE
}

func (d *DevClient) setToken(t string) {
	d.DevToken = t
}
func (d *DevClient) getToken() string {
	return d.DevToken
}
