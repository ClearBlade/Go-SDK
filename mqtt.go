package GoSDK

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"runtime/debug"
	"time"

	mqttTypes "github.com/clearblade/mqtt_parsing"
	mqtt "github.com/clearblade/paho.mqtt.golang"
)

const (
	//Mqtt QOS 0
	QOS_AtMostOnce = iota
	//Mqtt QOS 1
	QOS_AtLeastOnce
	//Mqtt QOS 2
	QOS_PreciselyOnce
	PUBLISH_HTTP_PREAMBLE = "/api/v/1/message/"
	_NEW_MH_PREAMBLE      = "/api/v/4/message/"
)

//LastWillPacket is a type to represent the Last Will and Testament packet
type LastWillPacket struct {
	Topic  string
	Body   string
	Qos    int
	Retain bool
}

type Callbacks struct {
	OnConnectCallback        mqtt.OnConnectHandler
	OnConnectionLostCallback mqtt.ConnectionLostHandler
}

func (b *client) NewClientID() string {
	buf := make([]byte, 10)
	rand.Read(buf)
	return fmt.Sprintf("%X", buf)
}

//herein we use the same trick we used for http clients

//InitializeMQTT allocates the mqtt client for the user. an empty string can be passed as the second argument for the user client
func (u *UserClient) InitializeMQTT(clientid string, ignore string, timeout int, ssl *tls.Config, lastWill *LastWillPacket) error {
	fmt.Printf("InitializeMQTT: %s\n", string(debug.Stack()))
	mqc, err := newMqttClient(u.UserToken, u.SystemKey, u.SystemSecret, clientid, timeout, u.MqttAddr, ssl, lastWill)
	if err != nil {
		return err
	}
	u.MQTTClient = mqc
	return nil
}

func (u *UserClient) InitializeMQTTWithCallback(clientid string, ignore string, timeout int, ssl *tls.Config, lastWill *LastWillPacket, callbacks *Callbacks) error {
	mqc, err := newMqttClientWithCallbacks(u.UserToken, u.SystemKey, u.SystemSecret, clientid, timeout, u.MqttAddr, ssl, lastWill, callbacks)
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
func (d *DevClient) InitializeMQTT(clientid, systemkey string, timeout int, ssl *tls.Config, lastWill *LastWillPacket) error {
	mqc, err := newMqttClient(d.DevToken, systemkey, "", clientid, timeout, d.MqttAddr, ssl, lastWill)
	if err != nil {
		return err
	}
	d.MQTTClient = mqc
	return nil
}

func (d *DevClient) InitializeMQTTWithCallback(clientid, systemkey string, timeout int, ssl *tls.Config, lastWill *LastWillPacket, callbacks *Callbacks) error {
	mqc, err := newMqttClientWithCallbacks(d.DevToken, systemkey, "", clientid, timeout, d.MqttAddr, ssl, lastWill, callbacks)
	if err != nil {
		return err
	}
	d.MQTTClient = mqc
	return nil
}

//InitializeMQTT allocates the mqtt client for the user. an empty string can be passed as the second argument for the user client
func (d *DeviceClient) InitializeMQTT(clientid string, ignore string, timeout int, ssl *tls.Config, lastWill *LastWillPacket) error {
	mqc, err := newMqttClient(d.DeviceToken, d.SystemKey, d.SystemSecret, clientid, timeout, d.MqttAddr, ssl, lastWill)
	if err != nil {
		return err
	}
	d.MQTTClient = mqc
	return nil
}

func (d *DeviceClient) InitializeMQTTWithCallback(clientid string, ignore string, timeout int, ssl *tls.Config, lastWill *LastWillPacket, callbacks *Callbacks) error {
	mqc, err := newMqttClientWithCallbacks(d.DeviceToken, d.SystemKey, d.SystemSecret, clientid, timeout, d.MqttAddr, ssl, lastWill, callbacks)
	if err != nil {
		return err
	}
	d.MQTTClient = mqc
	return nil
}

//Publish publishes a message to the specified mqtt topic
func (u *UserClient) Publish(topic string, message []byte, qos int) error {
	return publish(u.MQTTClient, topic, message, qos, u.getMessageId())
}

//Publish publishes a message to the specified mqtt topic
func (d *DeviceClient) Publish(topic string, message []byte, qos int) error {
	return publish(d.MQTTClient, topic, message, qos, d.getMessageId())
}

//Publish publishes a message to the specified mqtt topic
func (d *DevClient) Publish(topic string, message []byte, qos int) error {
	return publish(d.MQTTClient, topic, message, qos, d.getMessageId())
}

func (d *DevClient) PublishHttp(systemKey, topic string, message []byte, qos int) error {
	creds, err := d.credentials()
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"topic": topic,
		"body":  string(message[:]),
		"qos":   qos,
	}
	_, err = post(d, PUBLISH_HTTP_PREAMBLE+systemKey+"/publish", data, creds, nil)
	if err != nil {
		return err
	}
	return nil
}

//Subscribe subscribes a user to a topic. Incoming messages will be sent over the channel.
func (u *UserClient) Subscribe(topic string, qos int) (<-chan *mqttTypes.Publish, error) {
	return subscribe(u.MQTTClient, topic, qos)
}

//Subscribe subscribes a device to a topic. Incoming messages will be sent over the channel.
func (d *DeviceClient) Subscribe(topic string, qos int) (<-chan *mqttTypes.Publish, error) {
	return subscribe(d.MQTTClient, topic, qos)
}

//Subscribe subscribes a user to a topic. Incoming messages will be sent over the channel.
func (d *DevClient) Subscribe(topic string, qos int) (<-chan *mqttTypes.Publish, error) {
	return subscribe(d.MQTTClient, topic, qos)
}

//Unsubscribe stops the flow of messages over the corresponding subscription chan
func (u *UserClient) Unsubscribe(topic string) error {
	return unsubscribe(u.MQTTClient, topic)
}

//Unsubscribe stops the flow of messages over the corresponding subscription chan
func (d *DeviceClient) Unsubscribe(topic string) error {
	return unsubscribe(d.MQTTClient, topic)
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
func (d *DeviceClient) Disconnect() error {
	return disconnect(d.MQTTClient)
}

//Disconnect stops the TCP connection and unsubscribes the client from any remaining topics
func (d *DevClient) Disconnect() error {
	return disconnect(d.MQTTClient)
}

func (u *UserClient) GetCurrentTopicsWithQuery(systemKey string, columns []string, pageSize, pageNum int, descending bool) ([]map[string]interface{}, error) {
	return getMqttTopicsWithQuery(u, systemKey, columns, pageSize, pageNum, descending)
}

func (d *DevClient) GetCurrentTopicsWithQuery(systemKey string, columns []string, pageSize, pageNum int, descending bool) ([]map[string]interface{}, error) {
	return getMqttTopicsWithQuery(d, systemKey, columns, pageSize, pageNum, descending)
}

func (u *UserClient) GetCurrentTopicsCount(systemKey string) (map[string]interface{}, error) {
	return getMqttTopicsCount(u, systemKey)
}

func (d *DevClient) GetCurrentTopicsCount(systemKey string) (map[string]interface{}, error) {
	return getMqttTopicsCount(d, systemKey)
}

func (u *UserClient) GetCurrentTopics(systemKey string) ([]string, error) {
	return getMqttTopics(u, systemKey)
}

func (d *DevClient) GetCurrentTopics(systemKey string) ([]string, error) {
	return getMqttTopics(d, systemKey)
}

//Below are a series of convience functions to allow the user to only need to import
//the clearblade go-sdk
type mqttBaseClient struct {
	mqtt.Client
	address                                  string
	token, systemKey, systemSecret, clientID string
	timeout                                  int
}

//InitializeMqttClient allocates a mqtt client.
//the values for initialization are drawn from the client struct
//with the exception of the timeout and client id, which is mqtt specific.
// timeout refers to broker connect timeout
func newMqttClient(token, systemkey, systemsecret, clientid string, timeout int, address string, ssl *tls.Config, lastWill *LastWillPacket) (MqttClient, error) {
	o := mqtt.NewClientOptions()
	o.SetAutoReconnect(true)
	if ssl != nil {
		o.AddBroker("tls://" + address)
		o.SetTLSConfig(ssl)
	} else {
		o.AddBroker("tcp://" + address)
	}
	o.SetClientID(clientid)
	o.SetUsername(token)
	o.SetPassword(systemkey)
	o.SetConnectTimeout(time.Duration(timeout) * time.Second)
	if lastWill != nil {
		o.SetWill(lastWill.Topic, lastWill.Body, uint8(lastWill.Qos), lastWill.Retain)
	}
	cli := mqtt.NewClient(o)
	mqc := &mqttBaseClient{cli, address, token, systemkey, systemsecret, clientid, timeout}
	ret := mqc.Connect()
	ret.Wait()
	return mqc, ret.Error()
}

func newMqttClientWithCallbacks(token, systemkey, systemsecret, clientid string, timeout int, address string, ssl *tls.Config, lastWill *LastWillPacket, callbacks *Callbacks) (MqttClient, error) {
	o := mqtt.NewClientOptions()
	o.SetAutoReconnect(true)
	if ssl != nil {
		o.AddBroker("tls://" + address)
		o.SetTLSConfig(ssl)
	} else {
		o.AddBroker("tcp://" + address)
	}
	o.SetClientID(clientid)
	o.SetUsername(token)
	o.SetPassword(systemkey)
	o.SetConnectTimeout(time.Duration(timeout) * time.Second)
	if lastWill != nil {
		o.SetWill(lastWill.Topic, lastWill.Body, uint8(lastWill.Qos), lastWill.Retain)
	}
	if callbacks.OnConnectionLostCallback != nil {
		o.SetConnectionLostHandler(callbacks.OnConnectionLostCallback)
	}
	if callbacks.OnConnectCallback != nil {
		o.SetOnConnectHandler(callbacks.OnConnectCallback)
	}
	cli := mqtt.NewClient(o)
	mqc := &mqttBaseClient{cli, address, token, systemkey, systemsecret, clientid, timeout}
	ret := mqc.Connect()
	ret.Wait()
	return mqc, ret.Error()
}

func publish(c MqttClient, topic string, data []byte, qos int, mid uint16) error {
	if c == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	ret := c.Publish(topic, uint8(qos), false, data)
	return ret.Error()
}

func subscribe(c MqttClient, topic string, qos int) (<-chan *mqttTypes.Publish, error) {
	if c == nil {
		return nil, errors.New("MQTTClient is uninitialized")
	}
	pubs := make(chan *mqttTypes.Publish, 50)
	ret := c.Subscribe(topic, uint8(qos), func(client mqtt.Client, msg mqtt.Message) {
		path, _ := mqttTypes.NewTopicPath(msg.Topic())
		pubs <- &mqttTypes.Publish{Topic: path, Payload: msg.Payload()}
	})
	ret.WaitTimeout(1 * time.Second)
	return pubs, ret.Error()
}

func unsubscribe(c MqttClient, topic string) error {
	if c == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	ret := c.Unsubscribe(topic)
	ret.WaitTimeout(1 * time.Second)
	return ret.Error()
}

func disconnect(c MqttClient) error {
	if c == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	c.Disconnect(250)
	return nil
}

func getMqttTopicsWithQuery(c cbClient, systemKey string, columns []string, pageSize, pageNum int, descending bool) ([]map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}

	tmpQ := &Query{
		Columns:    columns,
		PageSize:   pageSize,
		PageNumber: pageNum,
		Order: []Ordering{
			Ordering{
				SortOrder: descending,
				OrderKey:  "topicid",
			},
		},
	}
	qry, err := createQueryMap(tmpQ)
	if err != nil {
		return nil, err
	}

	resp, err := get(c, _NEW_MH_PREAMBLE+systemKey+"/topics", qry, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}

	results, err := convertToMapStringInterface(resp.Body.([]interface{}))
	if err != nil {
		return nil, err
	}
	return results, nil
}

func getMqttTopicsCount(c cbClient, systemKey string) (map[string]interface{}, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}
	resp, err := get(c, _NEW_MH_PREAMBLE+systemKey+"/topics/count", nil, creds, nil)
	resp, err = mapResponse(resp, err)
	if err != nil {
		return nil, err
	}
	result, ok := resp.Body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("getMqttTopicsCount returns %T, expecting a map", resp.Body)
	}
	return result, nil
}

func convertToMapStringInterface(thing []interface{}) ([]map[string]interface{}, error) {
	rval := make([]map[string]interface{}, len(thing))
	for idx, vIF := range thing {
		switch vIF.(type) {
		case map[string]interface{}:
			rval[idx] = vIF.(map[string]interface{})
		default:
			return nil, fmt.Errorf("Bad type returned. Expecting a map, got %T", vIF)
		}
	}
	return rval, nil
}

func getMqttTopics(c cbClient, systemKey string) ([]string, error) {
	creds, err := c.credentials()
	if err != nil {
		return nil, err
	}

	resp, err := get(c, _MH_PREAMBLE+systemKey+"/currentTopics", nil, creds, nil)
	if err != nil {
		return nil, err
	}

	//parse the contents of the response body and return the topics in an array
	//Convert the array of interfaces to an array of strings
	topics := make([]string, len(resp.Body.([]interface{})))

	for i, topic := range resp.Body.([]interface{}) {
		topics[i] = topic.(string)
	}

	return topics, err
}
