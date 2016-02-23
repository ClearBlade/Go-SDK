package GoSDK

import (
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	mqtt "github.com/clearblade/mqtt_parsing"
	mqcli "github.com/clearblade/mqttclient"
	"time"
)

var (
	NetworkByteOrder = binary.BigEndian
)

const (
	//Mqtt QOS 0
	QOS_AtMostOnce = iota
	//Mqtt QOS 1
	QOS_AtLeastOnce
	//Mqtt QOS 2
	QOS_PreciselyOnce
)

//LastWillPacket is a type to represent the Last Will and Testament packet
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

//ConnectMQTT allows the user to connect to the mqtt broker. If no TLS config is provided, a TCP socket will be used
func (u *UserClient) ConnectMQTT(ssl *tls.Config, lastWill *LastWillPacket) error {
	//a questionable pointer, mainly for ease of checking nil
	//be more efficient to pass on the stack
	return connectToBroker(u.MQTTClient, CB_MSG_ADDR, ssl, lastWill)
}

//ConnectMQTT allows the user to connect to the mqtt broker. If no TLS config is provided, a TCP socket will be used
func (d *DevClient) ConnectMQTT(ssl *tls.Config, lastWill *LastWillPacket) error {
	return connectToBroker(d.MQTTClient, CB_MSG_ADDR, ssl, lastWill)
}

//Publish publishes a message to the specified mqtt topic
func (u *UserClient) Publish(topic string, message []byte, qos int) error {
	return publish(u.MQTTClient, topic, message, qos, u.getMessageId())
}

//Publish publishes a message to the specified mqtt topic
func (d *DevClient) Publish(topic string, message []byte, qos int) error {
	return publish(d.MQTTClient, topic, message, qos, d.getMessageId())
}

//Subscribe subscribes a user to a topic. Incoming messages will be sent over the channel.
func (u *UserClient) Subscribe(topic string, qos int) (<-chan *mqtt.Publish, error) {
	return subscribe(u.MQTTClient, topic, qos)
}

//Subscribe subscribes a user to a topic. Incoming messages will be sent over the channel.
func (d *DevClient) Subscribe(topic string, qos int) (<-chan *mqtt.Publish, error) {
	return subscribe(d.MQTTClient, topic, qos)
}

//Unsubscribe stops the flow of messages over the corresponding subscription chan
func (u *UserClient) Unsubscribe(topic string) error {
	return unsubscribe(u.MQTTClient, topic)
}

//Unsubscribe stops the flow of messages over the corresponding subscription chan
func (d *DevClient) Unsubscribe(topic string) error {
	return unsubscribe(d.MQTTClient, topic)
}

//Disconnect stops the TCP connection and unsubscribes the client from any remaining topics
func (u *UserClient) Disconnect() error {
	return disconnect(u.MQTTClient)
}

//Disconnect stops the TCP connection and unsubscribes the client from any remaining topics
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

func (u *UserClient) AuthenticateUsingMqtt() error {
	tok, err := authUsingMqtt(u.SystemKey, u.SystemSecret, u.Email, u.Password)
	if err != nil {
		return err
	} else {
		u.UserToken = tok
		return nil
	}
}

func (d *DevClient) AuthenticateUsingMqtt(systemKey, systemSecret string) error {
	tok, err := authUsingMqtt(systemKey, systemSecret, d.Email, d.Password)
	if err != nil {
		return err
	} else {
		d.DevToken = tok
		return nil
	}
}

func authUsingMqtt(systemKey, systemSecret, email, password string) (string, error) {
	//jumbling up the return values
	cli := mqcli.NewClient(systemKey, systemSecret, "", formatClientId(email, password), 30)
	err := cli.Start(CB_MSG_ADDR, &tls.Config{})
	if err != nil {
		return "", fmt.Errorf("error starting mqtt client:%s\n", err.Error())
	}
	err = mqcli.SendConnect(cli, false, false, 0, "", "")
	if err != nil {
		return "", fmt.Errorf("Error in connecting to broker for mqtt auth: %s\n", err)
	}
	// go ahead and get read to disconnect
	defer mqcli.SendDisconnect(cli)
	//we should be connected
	//now we subscribe
	mch, err := mqcli.SubscribeFlow(cli, formatTopicPath(systemKey, email), 0)
	if err != nil {
		return "", fmt.Errorf("Error in subscribing to broker for mqtt auth:%s\n", err.Error())
	}
	select {
	case msg := <-mch:
		token, _, actualServers, err := stripOutMqttAuthPayload(msg.Payload)
		if err != nil {
			return "", fmt.Errorf("poorly formatted mqtt authentication packet: %s\n", err.Error())
		}
		CB_MSG_ADDR = actualServers
		return token, err
	case <-cli.ClientErrorBuffer:
		return "", fmt.Errorf("bad return code. unauthorized\n")
	case <-time.After(time.Second * 20):
	}
	return "", fmt.Errorf("authenticating via mqtt timed out\n")
}

func formatClientId(u, p string) string {
	return u + ":" + p
}

func formatTopicPath(sk, u string) string {
	return sk + "/" + u
}

func stripOutMqttAuthPayload(payload []byte) (string, string, string, error) {
	siz_tok := NetworkByteOrder.Uint16(payload)
	if int(siz_tok) > len(payload)-2 {
		return "", "", "", fmt.Errorf("not enough space in payload")
	}
	tok := string(payload[2 : 2+siz_tok])
	siz_usrid := NetworkByteOrder.Uint16(payload[2+siz_tok:])
	if int(siz_tok+siz_usrid+4) > len(payload) {
		return "", "", "", fmt.Errorf("not enough space in payload")
	}
	usrid := string(payload[4+siz_tok : 4+siz_tok+siz_usrid])
	siz_url := NetworkByteOrder.Uint16(payload[4+siz_tok+siz_usrid:])
	if int(siz_tok+siz_usrid+siz_url+6) > len(payload) {
		return "", "", "", fmt.Errorf("not enough space in payload")
	}
	ip := string(payload[6+siz_tok+siz_usrid:])
	return tok, usrid, ip, nil
}
