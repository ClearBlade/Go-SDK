package GoSDK

import (
	"errors"
	"fmt"
	//	"log"
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
	resp, err := post(d.preamble()+"/systemmanagement", map[string]interface{}{
		"name":          name,
		"description":   description,
		"auth_required": users,
	}, creds, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating new system: %v", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error Creating new system: %v", resp.Body)
	}

	// TODO we need to make this json
	fmt.Printf("NEW SYSTEM RETURNS: %+v\n", resp.Body)
	return strings.TrimSpace(strings.Split(resp.Body.(string), ":")[1]), nil
}

func (d *DevClient) GetSystem(key string) (*System, error) {
	creds, err := d.credentials()
	if err != nil {
		return &System{}, err
	} else if len(creds) != 1 {
		return nil, fmt.Errorf("Error getting system: No DevToken Supplied")
	}
	sysResp, sysErr := get(d.preamble()+"/systemmanagement", map[string]string{"id": key}, creds, nil)
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
	resp, err := delete(d.preamble()+"/systemmanagement", map[string]string{"id": s}, creds, nil)
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
	resp, err := put(d.preamble()+"/systemmanagement", map[string]interface{}{
		"id":   system_key,
		"name": system_name,
	}, creds, nil)
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
	resp, err := put(d.preamble()+"/systemmanagement", map[string]interface{}{
		"id":          system_key,
		"description": system_description,
	}, creds, nil)
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
	resp, err := put(d.preamble()+"/systemmanagement", map[string]interface{}{
		"id":            system_key,
		"auth_required": true,
	}, creds, nil)
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
	resp, err := put(d.preamble()+"/systemmanagement", map[string]interface{}{
		"id":            system_key,
		"auth_required": false,
	}, creds, nil)
	if err != nil {
		return fmt.Errorf("Error changing system auth: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error changing system auth: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) DevUserInfo() (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d.preamble()+"/userinfo", nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting userdata: %v", err)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) NewConnectCollection(systemkey string, connectConfig connectCollection) (string, error) {
	creds, err := d.credentials()
	m := connectConfig.toMap()
	m["appID"] = systemkey
	m["name"] = connectConfig.tableName()
	if err != nil {
		return "", err
	}
	resp, err := post(d.preamble()+"/collectionmanagement", m, creds, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating collection: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error creating collection %v\n", resp.Body)
	}
	return resp.Body.(map[string]interface{})["collectionID"].(string), nil
}

func (d *DevClient) AlterConnectionDetails(systemkey string, connectConfig connectCollection) error {
	creds, err := d.credentials()
	out := make(map[string]interface{})
	m := connectConfig.toMap()
	out["appID"] = systemkey
	out["name"] = connectConfig.tableName()
	out["connectionStringMap"] = m
	resp, err := put(d.preamble()+"/collectionmanagement", out, creds, nil)
	if err != nil {
		return fmt.Errorf("Error creating collection: %s", err.Error())
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("Error creating collection %v\n", resp.Body)
	} else {
		return nil
	}
}

func (d *DevClient) NewCollection(systemKey, name string) (string, error) {
	creds, err := d.credentials()
	if err != nil {
		return "", err
	}
	resp, err := post(d.preamble()+"/collectionmanagement", map[string]interface{}{
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

func (d *DevClient) DeleteCollection(colId string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d.preamble()+"/collectionmanagement", map[string]string{
		"id": colId,
	}, creds, nil)
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
	resp, err := put(d.preamble()+"/collectionmanagement", map[string]interface{}{
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

func (d *DevClient) DeleteColumn(collection_id, column_name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put(d.preamble()+"/collectionmanagement", map[string]interface{}{
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

func (d *DevClient) GetCollectionInfo(collection_id string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return map[string]interface{}{}, err
	}
	resp, err := get(d.preamble()+"/collectionmanagement", map[string]string{
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

//get collections list in system
func (d *DevClient) GetAllCollections(SystemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d.preamble()+"/allcollections", map[string]string{
		"appid": SystemKey,
	}, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting collection info: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting collection info: %v", resp.Body)
	}

	//fmt.Printf("body: %+v\n", resp.Body)
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) GetAllRoles(SystemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d.preamble()+"/user/"+SystemKey+"/roles", map[string]string{
		"appid": SystemKey,
	}, creds, nil)

	return resp.Body.([]interface{}), nil
}

func (d *DevClient) CreateRole(systemKey, role_id string) (interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{
		"name":        role_id,
		"collections": []map[string]interface{}{},
		"topics":      []map[string]interface{}{},
		"services":    []map[string]interface{}{},
	}
	resp, err := post(d.preamble()+"/user/"+systemKey+"/roles", data, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error updating a role to have a collection: %v", resp.Body)
	}
	return resp.Body, nil
}

func (d *DevClient) DeleteRole(systemKey, roleId string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d.preamble()+"/user/"+systemKey+"/roles", map[string]string{"role": roleId}, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting role: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) DeleteUser(systemKey, userId string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d.preamble()+"/user/"+systemKey, map[string]string{"user": userId}, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting user: %v", resp.Body)
	}

	return nil
}

func (d *DevClient) AddUserToRoles(systemKey, userId string, roles []string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"user": userId,
		"changes": map[string]interface{}{
			"roles": map[string]interface{}{
				"add": roles,
			},
		},
	}
	resp, err := put(d.preamble()+"/user/"+systemKey, data, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error adding roles to a user: %v", resp.Body)
	}

	return nil
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
	resp, err := put(d.preamble()+"/user/"+systemKey+"/roles", data, creds, nil)
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
	resp, err := put(d.preamble()+"/user/"+systemKey+"/roles", data, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating a role to have a service: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) AddGenericPermissionToRole(systemKey, role_id, permission string, level int) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"id": role_id,
		"changes": map[string]interface{}{
			"services":    []map[string]interface{}{},
			"topics":      []map[string]interface{}{},
			"collections": []map[string]interface{}{},
		},
	}

	data["changes"].(map[string]interface{})[permission] = map[string]interface{}{
		"permissions": level,
	}

	resp, err := put(d.preamble()+"/user/"+systemKey+"/roles", data, creds, nil)
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
