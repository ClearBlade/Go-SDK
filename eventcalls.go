package GoSDK

import (
	"fmt"
)

const (
	_EVENTS_DEFS_PREAMBLE  = "/admin/triggers/definitions"
	_EVENTS_HDLRS_PREAMBLE = "/admin/triggers/handlers/"
	_TIMERS_HDLRS_PREAMBLE = "/admin/triggers/timers/"
	_MH_PREAMBLE           = "/api/v/1/message/"
)

//GetEventDefinitions returns a slice of the event definitions
func (d *DevClient) GetEventDefinitions() ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_EVENTS_DEFS_PREAMBLE, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

//GetEventHandlers returns a slice of the event handlers for a system
func (d *DevClient) GetEventHandlers(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_EVENTS_HDLRS_PREAMBLE+systemKey, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

//GetEventHandler reuturns a single event handler
func (d *DevClient) GetEventHandler(systemKey, name string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_EVENTS_HDLRS_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

//CreateEventHandler creates an event handler, otherwise known as a trigger
func (d *DevClient) CreateEventHandler(systemKey, name string,
	data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(_EVENTS_HDLRS_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

//DeleteEventHandler removes the event handler
func (d *DevClient) DeleteEventHandler(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(_EVENTS_HDLRS_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

//UpdateEventHandler allows the developer to alter the code executed by the event handler
func (d *DevClient) UpdateEventHandler(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(_EVENTS_HDLRS_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

func mapResponse(resp *CbResp, err error) (*CbResp, error) {
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(resp.Body.(string))
	}
	return resp, nil
}

////////////////////////////////////////////////////////////////////////////////
//
//  Timer calls are from here down

//Returns a slice of timer descriptions
func (d *DevClient) GetTimers(systemKey string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_TIMERS_HDLRS_PREAMBLE+systemKey, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.([]interface{}), nil
}

//GetTimer returns the definition of a single timer
func (d *DevClient) GetTimer(systemKey, name string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_TIMERS_HDLRS_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

//CreateTimer allows the user to create the timer with code
func (d *DevClient) CreateTimer(systemKey, name string,
	data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := post(_TIMERS_HDLRS_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

//DeleteTimer removes the timer
func (d *DevClient) DeleteTimer(systemKey, name string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	resp, err := delete(_TIMERS_HDLRS_PREAMBLE+systemKey+"/"+name, nil, creds, nil)
	_, err = mapResponse(resp, err)
	return err
}

//UpdateTimer allows the developer to change the code executed with the timer
func (d *DevClient) UpdateTimer(systemKey, name string, data map[string]interface{}) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := put(_TIMERS_HDLRS_PREAMBLE+systemKey+"/"+name, data, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp.Body.(map[string]interface{}), nil
}

//MessageHistory allows the developer to retrieve the message history
func (d *DevClient) MessageHistory(systemKey string) (interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(_MH_PREAMBLE+systemKey, nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
