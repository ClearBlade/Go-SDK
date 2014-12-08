package GoSDK

import (
	"crypto/tls"
	"errors"
	mqtt "github.com/clearblade/mqtt_parsing"
	mqcli "github.com/clearblade/mqttclient"
	"sync"
)

const (
	QOS_AtMostOnce = iota
	QOS_AtLeastOnce
	QOS_PreciselyOnce
)

//MqttMessage is a wrapper around the mqtt Publish packet type
//for simplicity in importing libraries when consuming the Go-SDK
type MqttMessage struct {
	Payload   []byte
	MessageId int
	Topic     string
}

//herein we use the same trick we used for http clients

//InitializeMQTT allocates the mqtt client for the user. an empty string can be passed as the second argument for the user client
func (u *UserClient) InitializeMQTT(clientid string, ignore string, timeout int) error {
	mqc, err := initializeMqttClient(u.UserToken, u.SystemSecret, u.SystemKey, clientid, timeout)
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

func (u *UserClient) ConnectMQTT(ssl *tls.Config) error {
	return connectToBroker(u.MQTTClient, CB_MSG_ADDR, ssl)
}

func (d *DevClient) ConnectMQTT(ssl *tls.Config) error {
	return connectToBroker(d.MQTTClient, CB_MSG_ADDR, ssl)
}

func (u *UserClient) Publish(topic string, message []byte, qos int) error {
	return publish(u.MQTTClient, topic, message, qos, u.getMessageId())
}

func (d *DevClient) Publish(topic string, message []byte, qos int) error {
	return publish(d.MQTTClient, topic, message, qos, d.getMessageId())
}

func (u *UserClient) Subscribe(topic string, qos int) (<-chan *MqttMessage, error) {
	ch, err := subscribe(u.MQTTClient, topic, qos)
	if err != nil {
		return nil, err
	} else {
		out, ech := convertMqttMessage(ch)
		u.addpipe(topic, ech)
		return out, nil
	}

}

func (d *DevClient) Subscribe(topic string, qos int) (<-chan *MqttMessage, error) {
	ch, err := subscribe(d.MQTTClient, topic, qos)
	if err != nil {
		return nil, err
	} else {
		out, ech := convertMqttMessage(ch)
		d.addpipe(topic, ech)
		return out, nil
	}
}

func (u *UserClient) Unsubscribe(topic string) error {
	u.killpipe(topic)
	return unsubscribe(u.MQTTClient, topic)
}

func (d *DevClient) Unsubscribe(topic string) error {
	d.killpipe(topic)
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
func initializeMqttClient(token, username, password, clientid string, timeout int) (*mqcli.Client, error) {
	cli := mqcli.NewClient(token,
		username,
		password,
		clientid,
		timeout)
	return cli, nil
}

//ConnectToBroker connects to the broker and sends the connect packet
func connectToBroker(c *mqcli.Client, address string, ssl *tls.Config) error {
	if c == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	err := c.Start(address, ssl)
	if err != nil {
		return err
	}
	return mqcli.SendConnect(c)
}

func publish(c *mqcli.Client, topic string, data []byte, qos int, mid uint16) error {
	if c == nil {
		return errors.New("MQTTClient is uninitialized")
	}
	pub := mqcli.MakeMeABytePublish(topic, data, mid)
	return mqcli.PublishFlow(c, pub)
}

//Subscribe is a simple wrapper around the mqtt client library
func subscribe(c *mqcli.Client, topic string, qos int) (<-chan mqtt.Message, error) {
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

//TODO:have a way of keeping of when to kill these goroutines
func convertMqttMessage(msg <-chan mqtt.Message) (<-chan *MqttMessage, chan<- struct{}) {
	mmc := make(chan *MqttMessage, len(msg))
	ech := make(chan struct{}, 1)
	go func(pc <-chan mqtt.Message, mc chan *MqttMessage, erch chan struct{}) {
		//TODO: can we do this without the overhead of a goroutine?
		for {
			select {
			case msg := <-pc:
				pub, _ := msg.(*mqtt.Publish)
				newmm := &MqttMessage{
					Payload:   pub.Payload,
					MessageId: int(pub.MessageId),
					Topic:     pub.Topic.Whole,
				}
				mc <- newmm
			case <-erch:
				//close both channels in unsubscribe calls
				return
			}
		}
	}(msg, mmc, ech)
	return mmc, ech
}

type pipedict struct {
	p      map[string]chan<- struct{}
	pd_mut *sync.Mutex
}

func newpipedict() *pipedict {
	return &pipedict{
		pd_mut: new(sync.Mutex),
		p:      make(map[string]chan<- struct{}),
	}
}

func (pd *pipedict) addPipe(top string, ch chan<- struct{}) {
	pd.pd_mut.Lock()
	pd.p[top] = ch
	pd.pd_mut.Unlock()
}

func (pd *pipedict) removePipe(top string) {
	pd.pd_mut.Lock()
	ech := pd.p[top]
	delete(pd.p, top)
	ech <- struct{}{}
	pd.pd_mut.Unlock()
	//do we close the other chan?
	close(ech)
}
