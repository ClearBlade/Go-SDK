package GoSDK

import (
	"crypto/tls"
	"errors"
	mqtt "github.com/clearblade/mqtt_parsing"
	mqcli "github.com/clearblade/mqttclient"
)

const (
	QOS_AtMostOnce = iota
	QOS_AtLeastOnce
	QOS_PreciselyOnce
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

//Below are a series of convience functions to allow the user to only need to import
//the clearblade go-sdk

//ConnectToBroker connects to the broker and sends the connect packet
func (c *Client) ConnectToBroker(address string, ssl *tls.Config) error {
	if c.MQTTClient == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	err := c.MQTTClient.Start(address, ssl)
	if err != nil {
		return err
	}
	return mqcli.SendConnect(c.MQTTClient)
}

//PublishString is a simple helper to create a publish with a string payload
func (c *Client) PublishString(topic, data string, qos int) error {
	if c.MQTTClient == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	pub := mqcli.MakeMeAPublish(topic, data, uint16(c.mrand.Int()))
	return mqcli.PublishFlow(c.MQTTClient, pub)
}

func (c *Client) PublishBytes(topic string, data []byte, qos int) error {
	if c.MQTTClient == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	pub := mqcli.MakeMeABytePublish(topic, data, uint16(c.mrand.Int()))
	return mqcli.PublishFlow(c.MQTTClient, pub)
}

//Subscribe is a simple wrapper around the mqtt client library
func (c *Client) Subscribe(topic string, qos int) (<-chan mqtt.Message, error) {
	if c.MQTTClient == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	return mqcli.SubscribeFlow(c.MQTTClient, topic, qos)
}

//Unsubscribe is a simple wrapper around the mqtt client library
func (c *Client) Unsubscribe(topic string) error {
	if c.MQTTClient == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	return mqcli.UnsubscribeFlow(c.MQTTClient, topic)
}

//Disconnect is a simple wrapper for sending mqtt disconnects
func (c *Client) Disconnect() error {
	if c.MQTTClient == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	return mqcli.SendDisconnect(c.MQTTClient)
}
