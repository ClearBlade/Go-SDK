package GoSDK

import (
	"fmt"
	"strings"
)

const (
	_CODE_ADMIN_PREAMBLE    = "/admin/code/v/1"
	_CODE_ADMIN_PREAMBLE_V2 = "/codeadmin/v/2"
)

//Service is a helper struct for grouping facts about a code service
type Service struct {
	Name    string
	Code    string
	Version int
	Params  []string
	System  string
}

//CodeLog provides structure to the code log return value
type CodeLog struct {
	Log  string
	Time string
}

//GetServiceNames retrieves the service names for a particular system
func (d *DevClient) GetServiceNames(systemKey string) ([]string, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_CODE_ADMIN_PREAMBLE+"/"+systemKey, nil, creds, nil)
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

//GetService returns information about a specified service
func (d *DevClient) GetService(systemKey, name string) (*Service, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_CODE_PREAMBLE+"/"+systemKey+"/"+name, nil, creds, nil)
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

//SetServiceEffectiveUser allows the developer to set the userid that a service executes under.
func (d *DevClient) SetServiceEffectiveUser(systemKey, name, userid string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := put(_CODE_ADMIN_PREAMBLE+"/"+systemKey+"/"+name, map[string]interface{}{
		"runuser": userid,
	}, creds, nil)
	if err != nil {
		return fmt.Errorf("Error updating service: %v\n", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating service: %v\n", resp.Body)
	}
	return nil
}

//UpdateService facillitates changes to the service's code
func (d *DevClient) UpdateService(systemKey, name, code string, params []string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	code = strings.Replace(code, "\\n", "\n", -1) // just to make sure we're not creating a \\\n since we could have removed some of the double escapes
	code = strings.Replace(code, "\n", "\\n", -1) // add back in the escaped stuff
	resp, err := put(_CODE_ADMIN_PREAMBLE+"/"+systemKey+"/"+name, map[string]interface{}{"code": code, "parameters": params, "name": name}, creds, nil)
	if err != nil {
		return fmt.Errorf("Error updating service: %v\n", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating service: %v\n", resp.Body)
	}
	return nil
}

//NewServiceWithLibraries creates a new service with the specified code, params, and libraries/dependencies.
//Parameters is a slice of strings of parameter names
func (d *DevClient) NewServiceWithLibraries(systemKey, name, code, deps string, params []string) error {
	extra := map[string]interface{}{"parameters": params, "dependencies": deps}
	return d.newService(systemKey, name, code, extra)
}

//NewService creates a new service with a new name, code and params
func (d *DevClient) NewService(systemKey, name, code string, params []string) error {
	extra := map[string]interface{}{"parameters": params}
	return d.newService(systemKey, name, code, extra)
}

//EnableLogsForService activates logging for execution of a service
func (d *DevClient) EnableLogsForService(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	_, err = post(_CODE_ADMIN_PREAMBLE_V2+"/logs/"+systemKey+"/"+name, map[string]interface{}{"logging": "true"}, creds, nil)
	return err
}

//DisableLogsForService turns logging off for that service
func (d *DevClient) DisableLogsForService(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	_, err = post(_CODE_ADMIN_PREAMBLE_V2+"/logs/"+systemKey+"/"+name, map[string]interface{}{"logging": false}, creds, nil)
	return err
}

//AreServiceLogsEnabled allows the developer to query the state of logging
func (d *DevClient) AreServiceLogsEnabled(systemKey, name string) (bool, error) {
	creds, err := d.credentials()
	if err != nil {
		return false, err
	}
	resp, err := get(_CODE_ADMIN_PREAMBLE_V2+"/logs/"+systemKey+"/"+name+"/active", nil, creds, nil)
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

//GetLogsForService retrieves the logs for the service
func (d *DevClient) GetLogsForService(systemKey, name string) ([]CodeLog, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_CODE_ADMIN_PREAMBLE_V2+"/logs/"+systemKey+"/"+name, nil, creds, nil)
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

func (d *DevClient) newService(systemKey, name, code string, extra map[string]interface{}) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	//code = strings.Replace(code, "\\n", "\n", -1)
	//code = strings.Replace(code, "\n", "\\n", -1)
	extra["code"] = code
	resp, err := post(_CODE_ADMIN_PREAMBLE+"/"+systemKey+"/"+name, extra, creds, nil)
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
	resp, err := delete(_CODE_ADMIN_PREAMBLE+"/"+systemKey+"/"+name, nil, creds, nil)
	if err != nil {
		return fmt.Errorf("Error deleting service: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting service: %v", resp.Body)
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
