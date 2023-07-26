package brokers

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fernandoglatz/home-management/backend/configs"
	"github.com/fernandoglatz/home-management/backend/utils"
)

func Setup() error {
	topic := configs.ApplicationConfig.Mosquitto.Topic
	return utils.SubscribeMQTT(topic, onMessageReceived)
}

func onMessageReceived(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	payloadStr := string(payload)

	log.Println("Received message [" + payloadStr + "]")
}
