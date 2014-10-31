package GoSDK

import (
	"errors"
	mqcli "github.com/clearblade/mqttclient"
)

//InitializeMqttClient allocates a mqtt client.
//the values for initialization are drawn from the client struct
//with the exception of the timeout and client id, which is mqtt specific.
func (c *Client) InitializeMqttClient(clientid string, timeout int) error {
	var tok string
	if ut := c.GetUserToken(); ut != "" {
		tok = ut
	} else if dt := c.GetDevToken(); dt != "" {
		tok = dt
	} else if c.SystemKey == "" && tok == "" {
		//if there is no token, we are presently unable to log in
		//there is a method of pure mqtt authentication in development
		//but presently we require a token to be acquired via http
		return errors.New("System Key nor Token not present, cannot initialize mqtt")
		//The client is presently (2014-10-31) allowed
		//to try to initialize with purely a system key and secret.
		//That is discontinued legacy behavior
		//so don't rely on it!
	}
	c.MQTTClient = mqcli.NewClient(tok,
		c.SystemKey,
		c.SystemSecret,
		clientid,
		timeout)
	return nil
}
