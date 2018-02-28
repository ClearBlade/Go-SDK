package GoSDK

import (
	"fmt"
)

const (
	_CODE_PREAMBLE      = "/api/v/1/code"
	_CODE_USER_PREAMBLE = "/api/v/3/code"
)

func callService(c cbClient, systemKey, name string, params map[string]interface{}, log bool) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	var resp *CbResp
	if log {

		resp, err = post(c, _CODE_PREAMBLE+"/"+systemKey+"/"+name, params, creds, map[string][]string{"Logging-enabled": []string{"true"}})
	} else {
		resp, err = post(c, _CODE_PREAMBLE+"/"+systemKey+"/"+name, params, creds, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("Error calling %s service: %v", name, err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error calling %s service: %v", name, resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func createService(c cbClient, systemKey, name, code string, extra map[string]interface{}) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	extra["code"] = code
	resp, err := post(c, _CODE_USER_PREAMBLE+"/"+systemKey+"/service/"+name, extra, creds, nil)
	if err != nil {
		return fmt.Errorf("Error creating new service: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error creating new service: %v", resp.Body)
	}
	return nil
}

func deleteService(c cbClient, systemKey, name string) error {
	creds, err := c.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(c, _CODE_USER_PREAMBLE+"/"+systemKey+"/service/"+name, nil, creds, nil)
	if err != nil {
		return fmt.Errorf("Error deleting service: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error deleting service: %v", resp.Body)
	}
	return nil
}

func updateService(c cbClient, sysKey, name, code string, extra map[string]interface{}) (error, map[string]interface{}) {
	creds, err := c.credentials()
	if err != nil {
		return err, nil
	}
	resp, err := put(c, _CODE_USER_PREAMBLE+"/"+sysKey+"/service/"+name, extra, creds, nil)
	if err != nil {
		return fmt.Errorf("Error updating service: %v\n", err), nil
	}
	body, ok := resp.Body.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Service not created. First create service..."), nil
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error updating service: %v\n", resp.Body), nil
	}
	return nil, body
}

// func (u *UserClient) CreateEventHandler(systemKey, name string,
// 	data map[string]interface{}) (map[string]interface{}, error) {
// 	creds, err := u.credentials()
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, err := mapResponse(post(u, _CODE_USER_PREAMBLE+"/"+systemKey+"/trigger/"+name, data, creds, nil))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.Body.(map[string]interface{}), nil
// }

//DeleteEventHandler removes the event handler
func (u *UserClient) DeleteEventHandler(systemKey, name string) error {
	creds, err := u.credentials()
	if err != nil {
		return err
	}
	_, err = mapResponse(delete(u, _CODE_USER_PREAMBLE+"/"+systemKey+"/trigger/"+name, nil, creds, nil))
	return err
}

func (u *UserClient) UpdateEventHandler(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := (put(u, _CODE_USER_PREAMBLE+"/"+systemKey+"/trigger/"+name, data, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

// func (u *UserClient) GetEventHandler(systemKey, name string) (map[string]interface{}, error) {
// 	creds, err := u.credentials()
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, err := mapResponse(get(u, _CODE_USER_PREAMBLE+"/"+systemKey+"/trigger/"+name, nil, creds, nil))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.Body.(map[string]interface{}), nil
// }

func (u *UserClient) GetTimer(systemKey, name string) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := mapResponse(get(u, _CODE_USER_PREAMBLE+"/"+systemKey+"/timer/"+name, nil, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

//CreateTimer allows the user to create the timer with code
//Returns a single instance of the object described in GetTimers for the newly created timer
func (u *UserClient) CreateTimer(systemKey, name string,
	data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := mapResponse(post(u, _CODE_USER_PREAMBLE+"/"+systemKey+"/timer/"+name, data, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

//DeleteTimer removes the timer
func (u *UserClient) DeleteTimer(systemKey, name string) error {
	creds, err := u.credentials()
	if err != nil {
		return err
	}
	_, err = mapResponse(delete(u, _CODE_USER_PREAMBLE+"/"+systemKey+"/timer/"+name, nil, creds, nil))
	return err
}

//UpdateTimer allows the developer to change the code executed with the timer
//Returns an updated version of the timer as described in GetTimer
func (u *UserClient) UpdateTimer(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := u.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := mapResponse(put(u, _CODE_USER_PREAMBLE+"/"+systemKey+"/timer/"+name, data, creds, nil))
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

//CallService performs a call against the specific service with the specified parameters. The logging argument will allow the developer to call the service with logging enabled for just that run.
//The return value is a map[string]interface{} where the results will be stored in the key "results". If logs were enabled, they'll be in "log".
func (d *DevClient) CallService(systemKey, name string, params map[string]interface{}, log bool) (map[string]interface{}, error) {
	return callService(d, systemKey, name, params, log)
}

//CallService performs a call against the specific service with the specified parameters.
//The return value is a map[string]interface{} where the results will be stored in the key "results". If logs were enabled, they'll be in "log".
func (u *UserClient) CallService(systemKey, name string, params map[string]interface{}) (map[string]interface{}, error) {
	return callService(u, systemKey, name, params, false)
}

func (u *UserClient) CreateService(systemKey, name, code string, params []string) error {
	extra := map[string]interface{}{"parameters": params}
	return createService(u, systemKey, name, code, extra)
}

func (u *UserClient) DeleteService(systemKey, name string) error {
	return deleteService(u, systemKey, name)
}

func (u *UserClient) UpdateService(systemKey, name, code string, params []string) (error, map[string]interface{}) {
	extra := map[string]interface{}{"code": code, "name": name, "parameters": params}
	return updateService(u, systemKey, name, code, extra)
}

func (u *UserClient) CreateTrigger(systemKey, name string,
	data map[string]interface{}) (map[string]interface{}, error) {
	return u.CreateEventHandler(systemKey, name, data)
}

func (u *UserClient) DeleteTrigger(systemKey, name string) error {
	return u.DeleteEventHandler(systemKey, name)
}

func (u *UserClient) UpdateTrigger(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	return u.UpdateEventHandler(systemKey, name, data)
}

func (u *UserClient) GetTrigger(systemKey, name string) (map[string]interface{}, error) {
	return u.GetEventHandler(systemKey, name)
}
