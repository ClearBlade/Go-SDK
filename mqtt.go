package GoSDK

import (
	"crypto/tls"
	"errors"
	mqtt "github.com/clearblade/mqtt_parsing"
	mqcli "github.com/clearblade/mqttclient"
)

const (
	//Mqtt QOS 0
	QOS_AtMostOnce = iota
	//Mqtt QOS 1
	QOS_AtLeastOnce
	//Mqtt QOS 2
	QOS_PreciselyOnce
)

type LastWillPacket struct {
	Topic  string
	Body   string
	Qos    int
	Retain bool
}

//herein we use the same trick we used for http clients

//InitializeMQTT allocates the mqtt client for the user. an empty string can be passed as the second argument for the user client
func (u *UserClient) InitializeMQTT(clientid string, ignore string, timeout int) error {
	mqc, err := initializeMqttClient(u.UserToken, u.SystemKey, u.SystemSecret, clientid, timeout)
	if err != nil {
		return err
	}
	u.MQTTClient = mqc
	return nil
}

//InitializeMQTT allocates the mqtt client for the developer. the second argument is a
//the systemkey you wish to use for authenticating with the message broker
//topics are isolated across systems, so in order to communicate with a specific
//system, you must supply the system key
func (d *DevClient) InitializeMQTT(clientid, systemkey string, timeout int) error {
	mqc, err := initializeMqttClient(d.DevToken, systemkey, "", clientid, timeout)
	if err != nil {
		return err
	}
	d.MQTTClient = mqc
	return nil
}

func (u *UserClient) ConnectMQTT(ssl *tls.Config, lastWill *LastWillPacket) error {
	//a questionable pointer, mainly for ease of checking nil
	//be more efficient to pass on the stack
	return connectToBroker(u.MQTTClient, CB_MSG_ADDR, ssl, lastWill)
}

func (d *DevClient) ConnectMQTT(ssl *tls.Config, lastWill *LastWillPacket) error {
	return connectToBroker(d.MQTTClient, CB_MSG_ADDR, ssl, lastWill)
}

func (u *UserClient) Publish(topic string, message []byte, qos int) error {
	return publish(u.MQTTClient, topic, message, qos, u.getMessageId())
}

func (d *DevClient) Publish(topic string, message []byte, qos int) error {
	return publish(d.MQTTClient, topic, message, qos, d.getMessageId())
}

func (u *UserClient) Subscribe(topic string, qos int) (<-chan *mqtt.Publish, error) {
	return subscribe(u.MQTTClient, topic, qos)
}

func (d *DevClient) Subscribe(topic string, qos int) (<-chan *mqtt.Publish, error) {
	return subscribe(d.MQTTClient, topic, qos)
}

func (u *UserClient) Unsubscribe(topic string) error {
	return unsubscribe(u.MQTTClient, topic)
}

func (d *DevClient) Unsubscribe(topic string) error {
	return unsubscribe(d.MQTTClient, topic)
}

func (u *UserClient) Disconnect() error {
	return disconnect(u.MQTTClient)
}

func (d *DevClient) Disconnect() error {
	return disconnect(d.MQTTClient)
}

//Below are a series of convience functions to allow the user to only need to import
//the clearblade go-sdk

//InitializeMqttClient allocates a mqtt client.
//the values for initialization are drawn from the client struct
//with the exception of the timeout and client id, which is mqtt specific.
func initializeMqttClient(token, systemkey, systemsecret, clientid string, timeout int) (*mqcli.Client, error) {
	cli := mqcli.NewClient(token,
		systemkey,
		systemsecret,
		clientid,
		timeout)
	return cli, nil
}

//ConnectToBroker connects to the broker and sends the connect packet
func connectToBroker(c *mqcli.Client, address string, ssl *tls.Config, lastWill *LastWillPacket) error {
	if c == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	err := c.Start(address, ssl)
	if err != nil {
		return err
	}
	if lastWill == nil {
		return mqcli.SendConnect(c, false, false, 0, "", "")
	} else {
		return mqcli.SendConnect(c, true, lastWill.Retain, lastWill.Qos, lastWill.Topic, lastWill.Body)
	}
}

func publish(c *mqcli.Client, topic string, data []byte, qos int, mid uint16) error {
	if c == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	pub, err := mqcli.MakeMeABytePublish(topic, data, mid)
	pub.Header.QOS = uint8(qos)
	if err != nil {
		return err
	}
	return mqcli.PublishFlow(c, pub)
}

//Subscribe is a simple wrapper around the mqtt client library
func subscribe(c *mqcli.Client, topic string, qos int) (<-chan *mqtt.Publish, error) {
	if c == nil {
		return nil, errors.New("MQTTClient is uninitialized")
	}
	return mqcli.SubscribeFlow(c, topic, qos)
}

//Unsubscribe is a simple wrapper around the mqtt client library
func unsubscribe(c *mqcli.Client, topic string) error {
	if c == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	return mqcli.UnsubscribeFlow(c, topic)
}

//Disconnect is a simple wrapper for sending mqtt disconnects
func disconnect(c *mqcli.Client) error {
	if c == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	return mqcli.SendDisconnect(c)
}
