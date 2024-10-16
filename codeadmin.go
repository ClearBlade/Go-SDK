package GoSDK

import (
	"fmt"
	"strings"
)

const (
	_CODE_ADMIN_PREAMBLE    = "/admin/code/v/1"
	_CODE_ADMIN_PREAMBLE_V2 = "/codeadmin/v/2"
	_CODE_ADMIN_PREAMBLE_V3 = "/codeadmin/v/3"
)

// Service is a helper struct for grouping facts about a code service
type Service struct {
	Name    string
	Code    string
	Version int
	Params  []string
	System  string
}

// CodeLog provides structure to the code log return value
type CodeLog struct {
	Log  string `json:"log"`
	Time string `json:"service_execution_time"`
}

// GetServiceNames retrieves the service names for a particular system
func (d *DevClient) GetServiceNames(systemKey string) ([]string, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _CODE_ADMIN_PREAMBLE+"/"+systemKey, nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting services: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting services: %v", resp.Body)
	}
	code := resp.Body.(map[string]interface{})["code"]
	sliceBody, isSlice := code.([]interface{})
	if !isSlice && code != nil {
		return nil, fmt.Errorf("Error getting services: server returned unexpected response")
	}
	services := make([]string, len(sliceBody))
	for i, service := range sliceBody {
		services[i] = service.(string)
	}
	return services, nil
}

// GetService returns information about a specified service
func (d *DevClient) GetService(systemKey, name string) (*Service, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _CODE_PREAMBLE+"/"+systemKey+"/"+name, nil, creds, nil)
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

func (d *DevClient) GetServiceRaw(systemKey, name string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _CODE_PREAMBLE+"/"+systemKey+"/"+name, nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting service: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting service: %v", resp.Body)
	}
	mapBody := resp.Body.(map[string]interface{})
	/*
		paramsSlice := mapBody["params"].([]interface{})
		params := make([]string, len(paramsSlice))
		for i, param := range paramsSlice {
			params[i] = param.(string)
		}
	*/
	return mapBody, nil
}

// SetServiceEffectiveUser allows the developer to set the userid that a service executes under.
func (d *DevClient) SetServiceEffectiveUser(systemKey, name, userid string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put(d, _CODE_ADMIN_PREAMBLE+"/"+systemKey+"/"+name, map[string]interface{}{
		"run_user": userid,
	}, creds, nil)
	if err != nil {
		return fmt.Errorf("Error updating service: %v\n", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating service: %v\n", resp.Body)
	}
	return nil
}

// UpdateService facillitates changes to the service's code
func (d *DevClient) UpdateService(systemKey, name, code string, params []string) (map[string]interface{}, error) {
	extra := map[string]interface{}{"code": code, "name": name, "parameters": params}
	return d.updateService(systemKey, name, code, extra)
}

func (d *DevClient) UpdateServiceWithLibraries(systemKey, name, code, deps string, params []string) (map[string]interface{}, error) {
	extra := map[string]interface{}{"code": code, "name": name, "parameters": params, "dependencies": deps}
	return d.updateService(systemKey, name, code, extra)
}

// UpdateServiceWithBody updates a service code and any additional info passed in (parameters, execution_timeout, etc.)
func (d *DevClient) UpdateServiceWithBody(systemKey, name, code string, extra map[string]interface{}) (map[string]interface{}, error) {
	return d.updateService(systemKey, name, code, extra)
}

func (d *DevClient) updateService(sysKey, name, code string, extra map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	extra["code"] = code
	resp, err := put(d, _CODE_ADMIN_PREAMBLE+"/"+sysKey+"/"+name, extra, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error updating service: %v\n", err)
	}
	body, ok := resp.Body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Service not created. First create service...")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error updating service: %v\n", resp.Body)
	}
	return body, nil
}

// NewServiceWithLibraries creates a new service with the specified code, params, and libraries/dependencies.
// Parameters is a slice of strings of parameter names
func (d *DevClient) NewServiceWithLibraries(systemKey, name, code, deps string, params []string) error {
	extra := map[string]interface{}{"parameters": params, "dependencies": deps}
	return d.newService(systemKey, name, code, extra)
}

// NewService creates a new service with a new name, code and params
func (d *DevClient) NewService(systemKey, name, code string, params []string) error {
	extra := map[string]interface{}{"parameters": params}
	return d.newService(systemKey, name, code, extra)
}

// NewServiceWithBody creates a new service with a new name, code and any additional info passed in (parameters, execution_timeout, etc.)
func (d *DevClient) NewServiceWithBody(systemKey, name, code string, extra map[string]interface{}) error {
	return d.newService(systemKey, name, code, extra)
}

// EnableLogsForService activates logging for execution of a service
func (d *DevClient) EnableLogsForService(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	_, err = post(d, _CODE_ADMIN_PREAMBLE_V2+"/logs/"+systemKey+"/"+name, map[string]interface{}{"logging": "true"}, creds, nil)
	return err
}

func (d *DevClient) UpdateServiceWithMap(systemKey, name string, changes map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(d, _CODE_ADMIN_PREAMBLE+"/"+systemKey+"/"+name, changes, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Error updating service: %v\n", err)
	}
	body, ok := resp.Body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Service not created. First create service...")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error updating service: %v\n", resp.Body)
	}
	return body, nil
}

// DisableLogsForService turns logging off for that service
func (d *DevClient) DisableLogsForService(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	_, err = post(d, _CODE_ADMIN_PREAMBLE_V2+"/logs/"+systemKey+"/"+name, map[string]interface{}{"logging": false}, creds, nil)
	return err
}

// AreServiceLogsEnabled allows the developer to query the state of logging
func (d *DevClient) AreServiceLogsEnabled(systemKey, name string) (bool, error) {
	creds, err := d.credentials()
	if err != nil {
		return false, err
	}
	resp, err := get(d, _CODE_ADMIN_PREAMBLE_V2+"/logs/"+systemKey+"/"+name+"/active", nil, creds, nil)
	if err != nil {
		return false, err
	}
	le := resp.Body.(map[string]interface{})["logging_enabled"]
	if le == nil {
		return false, fmt.Errorf("Improperly formatted json response")
	} else {
		return strings.ToLower(le.(string)) == "true", nil
	}
}

// GetLogsForService retrieves the logs for the service
func (d *DevClient) GetLogsForService(systemKey, name string) ([]CodeLog, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, _CODE_ADMIN_PREAMBLE_V2+"/logs/"+systemKey+"/"+name, nil, creds, nil)
	if err != nil {
		return nil, err
	}
	switch resp.Body.(type) {
	case string:
		return nil, fmt.Errorf("%s", resp.Body.(string))
	case []interface{}:
		r := resp.Body.([]map[string]interface{})
		outgoing := make([]CodeLog, len(r))
		for idx, v := range r {
			cl := genCodeLog(v)
			outgoing[idx] = cl
		}
		return outgoing, nil
	case []map[string]interface{}:
		r := resp.Body.([]map[string]interface{})
		outgoing := make([]CodeLog, len(r))
		for idx, v := range r {
			cl := genCodeLog(v)
			outgoing[idx] = cl
		}
		return outgoing, nil
	default:
		return nil, fmt.Errorf("Bad Return Value\n")
	}
}

func (d *DevClient) GetLogsForServiceInstance(systemKey, serviceName, serviceId string) (*CodeLog, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/logs/%s/%s/%s", _CODE_ADMIN_PREAMBLE_V2, systemKey, serviceName, serviceId)
	resp, err := get(d, url, nil, creds, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received error getting logs: %v", resp.Body)
	}

	var log CodeLog
	if err := decodeMapToStruct(resp.Body, &log); err != nil {
		return nil, fmt.Errorf("could not decode logs: %w", err)
	}

	return &log, nil
}

func (d *DevClient) SetLogLevel(systemKey, serviceName, logLevel string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	body := map[string]interface{}{"log_level": logLevel}
	url := fmt.Sprintf("%s/logs/%s/%s", _CODE_ADMIN_PREAMBLE_V2, systemKey, serviceName)
	resp, err := put(d, url, body, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("could not set log level: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received error response setting logs: %v", resp.Body)
	}

	if mapResp, ok := resp.Body.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("error setting log level; expected map got %T", resp.Body)
	} else {
		return mapResp, nil
	}
}

func (d *DevClient) newService(systemKey, name, code string, extra map[string]interface{}) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	extra["code"] = code
	resp, err := post(d, _CODE_ADMIN_PREAMBLE+"/"+systemKey+"/"+name, extra, creds, nil)
	if err != nil {
		return fmt.Errorf("Error creating new service: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error creating new service: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) DeleteService(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(d, _CODE_ADMIN_PREAMBLE+"/"+systemKey+"/"+name, nil, creds, nil)
	if err != nil {
		return fmt.Errorf("Error deleting service: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting service: %v", resp.Body)
	}
	return nil
}

func (d *DevClient) GetFailedServices(systemKey string) ([]map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(d, "/codeadmin/failed/"+systemKey, nil, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not get failed services: %s", err)
	}
	body := resp.Body.(map[string]interface{})[systemKey].([]interface{})
	services := make([]map[string]interface{}, len(body))
	for i, b := range body {
		services[i] = b.(map[string]interface{})
	}
	return services, nil
}

func (d *DevClient) RetryFailedServices(systemKey string, ids []string) ([]string, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(d, "/codeadmin/failed/"+systemKey, map[string]interface{}{"id": ids}, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not retry failed service %s/%s: %s", systemKey, ids, err)
	}
	body := resp.Body.([]interface{})
	responses := make([]string, len(body))
	for i, b := range body {
		responses[i] = b.(string)
	}
	return responses, nil
}

func (d *DevClient) DeleteFailedServices(systemKey string, ids []string) ([]map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := deleteWithBody(d, "/codeadmin/failed/"+systemKey, map[string]interface{}{"id": ids}, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not delete failed services %s/%s: %s", systemKey, ids, err)
	}
	body := resp.Body.([]interface{})
	services := make([]map[string]interface{}, len(body))
	for i, b := range body {
		services[i] = b.(map[string]interface{})
	}
	return services, nil
}

func (d *DevClient) SetLongRunningServiceParams(systemKey, name string, autoRestart, autoBalance bool, concurrency int) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	params := map[string]interface{}{
		"execution_timeout": -1,
		"auto_restart":      autoRestart,
		"concurrency":       concurrency,
		"auto_balance":      autoBalance,
	}

	_, err = put(d, _CODE_ADMIN_PREAMBLE_V2+"/"+systemKey+"/"+name, params, creds, nil)
	return err
}

func (d *DevClient) StopAllServiceInstances(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	services, err := d.GetRunningServices(systemKey)
	if err != nil {
		return err
	}

	for service := range services {
		params := map[string]string{
			"id": service,
		}

		_, err = deleteWithBody(d, _CODE_ADMIN_PREAMBLE_V3+"/running/"+systemKey, params, creds, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func genCodeLog(m map[string]interface{}) CodeLog {
	cl := CodeLog{}
	if tim, ext := m["service_execution_time"]; ext {
		t := tim.(string)
		cl.Time = t
	}
	if logg, ext := m["log"]; ext {
		l := logg.(string)
		cl.Log = l
	}
	return cl
}
