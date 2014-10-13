package GoSDK

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type System struct {
	Key         string
	Secret      string
	Name        string
	Description string
	Users       bool
}

func (c *Client) DevReg(email, password, fname, lname, org string) error {
	resp, err := c.Post("/admin/reg", map[string]interface{}{
		"email":    email,
		"password": password,
		"fname":    fname,
		"lname":    lname,
		"org":      org,
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%v", resp.Body))
	}
	c.SetDevToken(resp.Body.(map[string]interface{})["dev_token"].(string))
	return nil
}

func (c *Client) DevAuth(email, password string) error {
	resp, err := c.Post("/admin/auth", map[string]interface{}{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%v", resp.Body))
	}
	c.SetDevToken(resp.Body.(map[string]interface{})["dev_token"].(string))
	return nil
}

func (c *Client) DevLogout() error {
	if _, ok := c.Headers["ClearBlade-DevToken"]; !ok {
		return errors.New("No dev token stored. You need to auth first")
	}
	resp, err := c.Post("/admin/logout", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%d: %v", resp.StatusCode, resp.Body))
	}
	c.RemoveHeader("ClearBlade-DevToken")
	return nil
}

func (c *Client) NewSystem(name, description string, users bool) (string, error) {
	resp, err := c.Post("/admin/systemmanagement", map[string]interface{}{
		"name":          name,
		"description":   description,
		"auth_required": users,
	})
	if err != nil {
		return "", fmt.Errorf("Error creating new system: %v\n", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error Creating new system: %v\n", resp.Body)
	}

	// TODO we need to make this json
	return strings.TrimSpace(strings.Split(resp.Body.(string), ":")[1]), nil
}

func (c *Client) GetSystem(key string) (*System, error) {
	if _, ok := c.Headers["ClearBlade-DevToken"]; !ok {
		return nil, fmt.Errorf("Error getting system: No DevToken Supplied")
	}
	sysResp, sysErr := c.Get("/admin/systemmanagement", map[string]string{"id": key})
	if sysErr != nil {
		return nil, fmt.Errorf("Error gathering system information: %v\n", sysErr)
	}
	if sysResp.StatusCode != 200 {
		return nil, fmt.Errorf("Error gathering system information: %v\n", sysResp.Body)
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

func (c *Client) DeleteSystem(s string) error {
	resp, err := c.Delete("/admin/systemmanagement", map[string]string{"id": s})
	if err != nil {
		return fmt.Errorf("Error deleting system: %v\n", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting system: %v\n", resp.Body)
	}
	return nil
}

func (c *Client) DevUserInfo() error {
	resp, err := c.Get("/admin/userinfo", nil)
	if err != nil {
		return fmt.Errorf("Error getting userdata: %v\n", err)
	}
	log.Printf("HERE IS THE BODY: %+v\n", resp)
	return nil
}

func (c *Client) NewCollection(systemKey, name string) (string, error) {
	resp, err := c.Post("/admin/collectionmanagement", map[string]interface{}{
		"name":  name,
		"appID": systemKey,
	})
	if err != nil {
		return "", fmt.Errorf("Error creating collection: %v\n", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error creating collection %v\n", resp.Body)
	}
	return resp.Body.(map[string]interface{})["collectionID"].(string), nil
}

func (c *Client) DeleteCollection(colId string) error {
	resp, err := c.Delete("/admin/collectionmanagement", map[string]string{
		"id": colId,
	})
	if err != nil {
		return fmt.Errorf("Error deleting collection %v\n", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting collection %v\n", resp.Body)
	}
	return nil
}

func (c *Client) AddColumn(collection_id, column_name, column_type string) error {
	resp, err := c.Put("/admin/collectionmanagement", map[string]interface{}{
		"id": collection_id,
		"addColumn": map[string]interface{}{
			"name": column_name,
			"type": column_type,
		},
	})
	if err != nil {
		return fmt.Errorf("Error adding column: %v\n", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error adding column: %v\n", resp.Body)
	}
	return nil
}

func (c *Client) DeleteColumn(collection_id, column_name string) error {
	resp, err := c.Put("/admin/collectionmanagement", map[string]interface{}{
		"id":           collection_id,
		"deleteColumn": column_name,
	})
	if err != nil {
		return fmt.Errorf("Error deleting column: %v\n", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting column: %v\n", resp.Body)
	}
	return nil
}

func (c *Client) SetCollectionName(collection_id, collection_name string) error {
	resp, err := c.Put("/admin/collectionmanagement", map[string]interface{}{
		"id":   collection_id,
		"name": collection_name,
	})
	if err != nil {
		return fmt.Errorf("Error changing collection name: %v\n", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error changing collection name: %v\n", resp.Body)
	}
	return nil
}
