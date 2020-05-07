package main

import (
	"idp/Adaptor/database"
	"idp/Adaptor/messaging"
	log "github.com/sirupsen/logrus"
	"fmt"
	"os"
	"context"

	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
	"sync"
	"time"
)

const (
	InfluxDBAddress  = "http://influxdb:8086"
	InfluxDBName     = "IDP"
	InfluxDBUser     = "grafana_user"
	InfluxDBPassword = "influxdbGrafanaIDP"

	MQTTHost     = "tcp://mqtt:1883"
	MQTTClient   = "adaptor_client"
	MQTTUsername = ""
	MQTTPassword = ""
	MQTTTopic    = "#"
	MQTTQOS      = 0
)

var Logging = os.Getenv("DEBUG_DATA_FLOW")
var dataChannel = make(chan mqtt.Message)

type Adaptor struct {
	InfluxDB database.InfluxDB
	MQTT     *messaging.MQTTConnection
}

func (a *Adaptor) handleMessages() {
	for {
		message := <-dataChannel
		topic := message.Topic()
		payload := message.Payload()

		var data map[string]interface{}

		if err := json.Unmarshal(payload, &data); err != nil {
			if Logging == "true" {
				log.Errorf("Could not parse message data.")
			}
			continue
		}

		var timestamp time.Time

		timestamp = time.Now()

		employer := strings.Split(topic, "/")[1]

		var measurements = make([]string, 0)
		var fields = make([]map[string]interface{}, 0)
		var tags = make(map[string]string)
		tags["employer"] = employer

		for key, value := range data {
			switch value.(type) {
			case int, float32, float64:
				measurement := fmt.Sprintf("%s", key)
				measurements = append(measurements, measurement)

				fields = append(fields, make(map[string]interface{}))
				fields[len(fields)-1]["value"] = value

				if Logging == "true" {
					log.Infof("%s.%s.%s %v", employer, key, value)
				}
			default:
				continue
			}
		}

		if err := a.InfluxDB.AddMeasurements(context.Background(), measurements, fields, tags, timestamp); err != nil {
			if Logging == "true" {
				log.Error("Could not save data to database: ", err.Error())
			}

			continue
		}
	}
}

func on_message(client mqtt.Client, message mqtt.Message) {
	if Logging == "true" {
		log.Info("Received a message by topic " + message.Topic())
	}

	dataChannel <- message
}

func NewAdaptor() *Adaptor {
	return &Adaptor{
		InfluxDB: database.NewInfluxDBConnector(InfluxDBAddress, InfluxDBName, InfluxDBUser, InfluxDBPassword),
		MQTT:     messaging.NewMQTTConnection(MQTTHost, MQTTClient, MQTTUsername, MQTTPassword, MQTTTopic, MQTTQOS),
	}
}

func (a *Adaptor) Start() {
	go a.MQTT.Client.Subscribe(a.MQTT.ClientConf.Topic, a.MQTT.ClientConf.Qos, on_message)
	go a.handleMessages()
}

func main() {
	wg := sync.WaitGroup{}

	adaptor := NewAdaptor()

	wg.Add(1)
	adaptor.Start()
	wg.Wait()

	log.Fatal("Nu trebuie aici")
}
