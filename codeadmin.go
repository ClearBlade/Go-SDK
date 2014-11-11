package GoSDK

import (
	"fmt"
	"strings"
)

const (
	_CODE_ADMIN_PREAMBLE = "/admin/code/v/1"
)

type Service struct {
	Name    string
	Code    string
	Version int
	Params  []string
	System  string
}

func (d *DevClient) GetServiceNames(systemKey string) ([]string, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_CODE_ADMIN_PREAMBLE+"/"+systemKey, nil, creds)
	if err != nil {
		return nil, fmt.Errorf("Error getting services: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting services: %v", resp.Body)
	}
	code := resp.Body.(map[string]interface{})["code"]
	sliceBody, isSlice := code.([]interface{})
	if !isSlice {
		return nil, fmt.Errorf("Error getting services: server returned unexpected response")
	}
	services := make([]string, len(sliceBody))
	for i, service := range sliceBody {
		services[i] = service.(string)
	}
	return services, nil
}

func (d *DevClient) GetService(systemKey, name string) (*Service, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_CODE_PREAMBLE+"/"+systemKey+"/"+name, nil, creds)
	if err != nil {
		return nil, fmt.Errorf("Error getting service: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting service: %v", resp.Body)
	}
	mapBody := resp.Body.(map[string]interface{})
	paramsSlice := mapBody["params"].([]interface{})
	params := make([]string, len(paramsSlice))
	for i, param := range paramsSlice {
		params[i] = param.(string)
	}
	svc := &Service{
		Name:    name,
		System:  systemKey,
		Code:    mapBody["code"].(string),
		Version: int(mapBody["current_version"].(float64)),
		Params:  params,
	}
	return svc, nil
}

func (d *DevClient) UpdateService(systemKey, name, code string, params []string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	code = strings.Replace(code, "\\n", "\n", -1) // just to make sure we're not creating a \\\n since we could have removed some of the double escapes
	code = strings.Replace(code, "\n", "\\n", -1) // add back in the escaped stuff
	resp, err := put(_CODE_ADMIN_PREAMBLE+"/"+systemKey+"/"+name, map[string]interface{}{"code": code, "parameters": params}, creds)
	if err != nil {
		return fmt.Errorf("Error updating service: %v\n", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating service: %v\n", resp.Body)
	}
	return nil
}
