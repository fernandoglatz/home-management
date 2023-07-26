package utils

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fernandoglatz/home-management/backend/configs"
)

var client mqtt.Client

func ConnectToMQTT() error {
	log.Println("Connecting to MQTT")

	uri := configs.ApplicationConfig.Mosquitto.URI

	opts := mqtt.NewClientOptions()
	opts.AddBroker(uri)
	client = mqtt.NewClient(opts)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Println("Connected to MQTT!")

	return nil
}

func SubscribeMQTT(topic string, callback mqtt.MessageHandler) error {
	log.Println("Subscribing to topic " + topic)

	token := client.Subscribe(topic, 0, callback)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Println("Subscribed to topic " + topic + "!")

	return nil
}
