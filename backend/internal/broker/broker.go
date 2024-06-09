package broker

import (
	"context"
	"encoding/json"
	"fernandoglatz/home-management/internal/core/common/utils"
	"fernandoglatz/home-management/internal/core/common/utils/log"
	"fernandoglatz/home-management/internal/core/entity"
	"fernandoglatz/home-management/internal/infrastructure/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Setup(ctx context.Context) error {
	mqttBroker := utils.MqttBroker
	rabbitMqBroker := utils.RabbitMqBroker

	topic := config.ApplicationConfig.Broker.Mqtt.Topics.Broadcast
	err := mqttBroker.Subscribe(ctx, topic, onMqttBroadCastMessageReceived)
	if err != nil {
		return err
	}

	err = rabbitMqBroker.CreateQueue(ctx, "teste")
	if err != nil {
		return err
	}

	err = rabbitMqBroker.CreateExchange(ctx, "teste-exchange")
	if err != nil {
		return err
	}

	user := entity.User{Name: "Fernando"}

	err = rabbitMqBroker.PublishQueue(ctx, "teste", user)
	if err != nil {
		return err
	}

	err = rabbitMqBroker.Subscribe(ctx, "teste", onRabbitMqMessageReceived)
	if err != nil {
		return err
	}

	return nil
}

func onRabbitMqMessageReceived(queue string, delivery amqp.Delivery) error {
	ctx := context.Background()
	jsonData, _ := json.Marshal(delivery)
	json := string(jsonData)

	log.Info(ctx).Msg("Mensagem da fila: " + json)

	return nil
}

func onMqttBroadCastMessageReceived(client mqtt.Client, msg mqtt.Message) {
	ctx := context.Background()
	payload := msg.Payload()
	payloadStr := string(payload)

	log.Info(ctx).Msg("Received MQTT message [" + payloadStr + "]")
}
