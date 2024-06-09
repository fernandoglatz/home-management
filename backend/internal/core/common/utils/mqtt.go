package utils

import (
	"context"
	"encoding/json"
	"fernandoglatz/home-management/internal/core/common/utils/log"
	"fernandoglatz/home-management/internal/infrastructure/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const MQTT_QOS = 1

var MqttBroker MqttBrokerType

type MqttBrokerType struct {
	Client mqtt.Client
}

func ConnectToMQTT(ctx context.Context) error {
	log.Info(ctx).Msg("Connecting to MQTT")

	config := config.ApplicationConfig.Broker.Mqtt
	uri := config.Uri
	clientId := config.ClientId
	user := config.User
	password := config.Password

	opts := mqtt.NewClientOptions()
	opts.AddBroker(uri)
	opts.SetClientID(clientId)
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetCleanSession(false)
	client := mqtt.NewClient(opts)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	MqttBroker = MqttBrokerType{
		Client: client,
	}

	log.Info(ctx).Msg("Connected to MQTT!")

	return nil
}

func (mqttBroker *MqttBrokerType) Subscribe(ctx context.Context, topic string, callback mqtt.MessageHandler) error {
	log.Info(ctx).Msg("Subscribing to topic [" + topic + "]")

	client := mqttBroker.Client
	token := client.Subscribe(topic, MQTT_QOS, callback)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Info(ctx).Msg("Subscribed to topic [" + topic + "]!")

	return nil
}

func (mqttBroker *MqttBrokerType) Publish(ctx context.Context, topic string, object any) error {
	jsonData, err := json.Marshal(object)
	if err != nil {
		return err
	}
	json := string(jsonData)

	client := mqttBroker.Client
	token := client.Publish(topic, MQTT_QOS, false, json)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	log.Info(ctx).Msg("Published message [" + json + "] to topic [" + topic + "]")

	return nil
}
