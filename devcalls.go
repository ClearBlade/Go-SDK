package GoSDK

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	_DEV_HEADER_KEY = "ClearBlade-DevToken"
	_DEV_PREAMBLE   = "/admin"
)

type System struct {
	Key         string
	Secret      string
	Name        string
	Description string
	Users       bool
}

func (d *DevClient) NewSystem(name, description string, users bool) (string, error) {
	creds, err := d.credentials()
	if err != nil {
		return "", err
	}
	resp, err := post("/admin/systemmanagement", map[string]interface{}{
		"name":          name,
		"description":   description,
		"auth_required": users,
	}, creds)
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
	creds, err := d.credentials()
	if err != nil {
		return &System{}, err
	} else if len(creds) != 1 {
		return nil, fmt.Errorf("Error getting system: No DevToken Supplied")
	}
	sysResp, sysErr := get("/admin/systemmanagement", map[string]string{"id": key}, creds)
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
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete("/admin/systemmanagement", map[string]string{"id": s}, creds)
	if err != nil {
		return fmt.Errorf("Error deleting system: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting system: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) SetSystemName(system_key, system_name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put("/admin/systemmanagement", map[string]interface{}{
		"id":   system_key,
		"name": system_name,
	}, creds)
	if err != nil {
		return fmt.Errorf("Error changing system name: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error changing system name: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) SetSystemDescription(system_key, system_description string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put("/admin/systemmanagement", map[string]interface{}{
		"id":          system_key,
		"description": system_description,
	}, creds)
	if err != nil {
		return fmt.Errorf("Error changing system description: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error changing system description: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) SetSystemAuthOn(system_key string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put("/admin/systemmanagement", map[string]interface{}{
		"id":            system_key,
		"auth_required": true,
	}, creds)
	if err != nil {
		return fmt.Errorf("Error changing system auth: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error changing system auth: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) SetSystemAuthOff(system_key string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put("/admin/systemmanagement", map[string]interface{}{
		"id":            system_key,
		"auth_required": false,
	}, creds)
	if err != nil {
		return fmt.Errorf("Error changing system auth: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error changing system auth: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) DevUserInfo() error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := get("/admin/userinfo", nil, creds)
	if err != nil {
		return fmt.Errorf("Error getting userdata: %v", err)
	}
	log.Printf("HERE IS THE BODY: %+v\n", resp)
	return nil
}

func (d *DevClient) NewCollection(systemKey, name string) (string, error) {
	creds, err := d.credentials()
	if err != nil {
		return "", err
	}
	resp, err := post("/admin/collectionmanagement", map[string]interface{}{
		"name":  name,
		"appID": systemKey,
	}, creds)
	if err != nil {
		return "", fmt.Errorf("Error creating collection: %v", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error creating collection %v", resp.Body)
	}
	return resp.Body.(map[string]interface{})["collectionID"].(string), nil
}

func (d *DevClient) DeleteCollection(colId string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete("/admin/collectionmanagement", map[string]string{
		"id": colId,
	}, creds)
	if err != nil {
		return fmt.Errorf("Error deleting collection %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting collection %v", resp.Body)
	}
	return nil
}

func (d *DevClient) AddColumn(collection_id, column_name, column_type string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put("/admin/collectionmanagement", map[string]interface{}{
		"id": collection_id,
		"addColumn": map[string]interface{}{
			"name": column_name,
			"type": column_type,
		},
	}, creds)
	if err != nil {
		return fmt.Errorf("Error adding column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error adding column: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) DeleteColumn(collection_id, column_name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put("/admin/collectionmanagement", map[string]interface{}{
		"id":           collection_id,
		"deleteColumn": column_name,
	}, creds)
	if err != nil {
		return fmt.Errorf("Error deleting column: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting column: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) GetCollectionInfo(collection_id string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return map[string]interface{}{}, err
	}
	resp, err := get("/admin/collectionmanagement", map[string]string{
		"id": collection_id,
	}, creds)
	if err != nil {
		return nil, fmt.Errorf("Error getting collection info: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting collection info: %v", resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) AddCollectionToRole(systemKey, collection_id, role_id string, level int) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"id": role_id,
		"changes": map[string]interface{}{
			"collections": []map[string]interface{}{
				map[string]interface{}{
					"itemInfo": map[string]interface{}{
						"id": collection_id,
					},
					"permissions": level,
				},
			},
			"topics":   []map[string]interface{}{},
			"services": []map[string]interface{}{},
		},
	}
	resp, err := put("/admin/user/"+systemKey+"/roles", data, creds)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating a role to have a collection: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) AddServiceToRole(systemKey, service, role_id string, level int) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"id": role_id,
		"changes": map[string]interface{}{
			"services": []map[string]interface{}{
				map[string]interface{}{
					"itemInfo": map[string]interface{}{
						"name": service,
					},
					"permissions": level,
				},
			},
			"topics":      []map[string]interface{}{},
			"collections": []map[string]interface{}{},
		},
	}
	resp, err := put("/admin/user/"+systemKey+"/roles", data, creds)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating a role to have a service: %v", resp.Body)
	}
	return nil
}

//second verse, same as the first, eh?
func (d *DevClient) credentials() ([][]string, error) {
	if d.DevToken != "" {
		return [][]string{
			[]string{
				_DEV_HEADER_KEY,
				d.DevToken,
			},
		}, nil
	} else {
		return [][]string{}, errors.New("No SystemSecret/SystemKey combo, or UserToken found")
	}
}

func (d *DevClient) preamble() string {
	return _DEV_PREAMBLE
}

func (d *DevClient) getSystemInfo() (string, string) {
	return "", ""
}

func (d *DevClient) setToken(t string) {
	d.DevToken = t
}
func (d *DevClient) getToken() string {
	return d.DevToken
}

func (d *DevClient) getMessageId() uint16 {
	return uint16(d.mrand.Int())
}
