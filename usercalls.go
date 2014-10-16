package GoSDK

import (
	"fmt"
)

func (c *Client) UserReg(email, password string) error {
	resp, err := c.Post("/api/v/1/user/reg", map[string]interface{}{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return fmt.Errorf("Error registering user: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error registering user: %v", resp.Body)
	}
	return nil
}

func (c *Client) UserAuth(email, password string) error {
	resp, err := c.Post("/api/v/1/user/auth", map[string]interface{}{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return fmt.Errorf("Error authenticating user: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error authenticating user: %v", resp.Body)
	}
	c.SetUserToken(resp.Body.(map[string]interface{})["user_token"].(string))
	return nil
}

func (c *Client) UserLogout() error {
	if _, ok := c.Headers["ClearBlade-UserToken"]; !ok {
		return fmt.Errorf("No user token stored. you need to auth to get one")
	}
	resp, err := c.Post("/api/v/1/user/logout", nil)
	if err != nil {
		return fmt.Errorf("Error logging user out: %v", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error logging user out: %v", resp.Body)
	}
	return nil
}
