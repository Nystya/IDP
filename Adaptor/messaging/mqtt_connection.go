package messaging

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"time"
)

type MQTTConf struct {
	ClientOptions *mqtt.ClientOptions
	Topic         string
	Qos           byte
}

type MQTTConnection struct {
	Client     mqtt.Client
	ClientConf *MQTTConf
}

func NewMQTTConnection(host, clientID, username, password, topic string, qos byte) *MQTTConnection {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(host)
	opts.SetClientID(clientID)

	if username != "" && password != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()

	if !token.WaitTimeout(30 * time.Second) {
		log.Fatal("Connection to broker timed out.")
	}

	if err := token.Error(); err != nil {
		log.Fatal("Error: " + err.Error())
	}

	conf := &MQTTConf{
		ClientOptions: opts,
		Topic:         topic,
		Qos:           qos,
	}

	return &MQTTConnection{
		Client:     client,
		ClientConf: conf,
	}
}
