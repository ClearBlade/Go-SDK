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
