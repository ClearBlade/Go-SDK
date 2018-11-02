package GoSDK

import (
	"fmt"
)

//
// DevClient must already be platform admin to use these endpoints
//

func (d *DevClient) PromoteDevToPlatformAdmin(email string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"email": email,
	}

	resp, err := post(d, d.preamble()+"/promotedev", data, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error promoting %s to admin: %v", email, resp.Body)
	}
	return nil
}

func (d *DevClient) DemoteDevFromPlatformAdmin(email string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"email": email,
	}

	resp, err := post(d, d.preamble()+"/demotedev", data, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error demoting %s to admin: %v", email, resp.Body)
	}
	return nil
}

func (d *DevClient) ResetDevelopersPassword(email, newPassword string) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"email":        email,
		"new_password": newPassword,
	}

	resp, err := post(d, d.preamble()+"/resetpassword", data, creds, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error resetting %s's password: %v", email, resp.Body)
	}
	return nil
}

func (d *DevClient) GetSystemAnalytics(systemKey string) (interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	analyticsEndpoint := d.preamble() + "/platform/system/" + systemKey
	resp, err := get(d, analyticsEndpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting %s's analytics: %v", systemKey, resp.Body)
	}
	return resp.Body, nil
}

func (d *DevClient) GetAllSystemsAnalytics(query string) ([]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	analyticsEndpoint := d.preamble() + "/platform/systems"
	q := make(map[string]string)
	if len(query) != 0 {
		q["query"] = query
	} else {
		q = nil
	}
	resp, err := get(d, analyticsEndpoint, q, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error getting analytics: %v", resp.Body)
	}
	return resp.Body.([]interface{}), nil
}

func (d *DevClient) DisableSystem(systemKey string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	disableSystemEndpoint := d.preamble() + "/platform/" + systemKey
	resp, err := delete(d, disableSystemEndpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error disabling system %s: %v", systemKey, resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}

func (d *DevClient) EnableSystem(systemKey string) (map[string]interface{}, error) {
	creds, err := d.credentials()
	if err != nil {
		return nil, err
	}

	disableSystemEndpoint := d.preamble() + "/platform/" + systemKey
	resp, err := post(d, disableSystemEndpoint, nil, creds, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error disabling system %s: %v", systemKey, resp.Body)
	}
	return resp.Body.(map[string]interface{}), nil
}
