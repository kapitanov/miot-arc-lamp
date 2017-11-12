package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mqttOpts   mqtt.ClientOptions
	mqttClient mqtt.Client
)

const (
	topicRequest = "/arc/request"
	topicPush    = "/arc/state"

	mqttPort = 1883
)

func init() {
	mqttAddr, err := url.Parse(fmt.Sprintf("tcp://%s:%d", os.Getenv("MQTT_HOSTNAME"), mqttPort))
	if err != nil {
		panic(err)
	}
	mqttUsername := os.Getenv("MQTT_USERNAME")
	mqttPassword := os.Getenv("MQTT_PASSWORD")

	mqttOpts = mqtt.ClientOptions{
		Servers:  []*url.URL{mqttAddr},
		ClientID: "moit-arc-lamp",
		Username: mqttUsername,
		Password: mqttPassword,
	}
}

func runMqtt(reqUpdate chan int) error {
	mqttClient = mqtt.NewClient(&mqttOpts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		err := token.Error()
		fmt.Fprintf(os.Stderr, "mqtt: failed to connect. %s\n", err)
		return err
	}

	fmt.Fprintf(os.Stdout, "mqtt: connected\n")
	mqttClient.Subscribe(topicRequest, 0, func(_ mqtt.Client, msg mqtt.Message) {
		fmt.Fprintf(os.Stdout, "mqtt: recv. %s/%d\n", msg.Topic(), msg.MessageID())
		go func() { mqttPublish() }()
	})

	mqttPublish()

	return nil
}

func mqttPublish() {
	if mqttClient == nil {
		return
	}

	bytes, err := json.Marshal(currentStatus)
	if err == nil {
		token := mqttClient.Publish(topicPush, 0, false, bytes)
		token.Wait()
		err = token.Error()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "mqtt: failed to publish. %s\n", err)
	}

	fmt.Fprintf(os.Stdout, "mqtt: sent msg to %s\n", topicPush)
}