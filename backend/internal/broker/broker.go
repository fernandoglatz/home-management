package broker

import (
	"context"
	"fernandoglatz/home-management/internal/core/common/utils"
	"fernandoglatz/home-management/internal/core/common/utils/log"
	"fernandoglatz/home-management/internal/infrastructure/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Setup(ctx context.Context) error {
	topic := config.ApplicationConfig.Broker.Mqtt.Topics.Broadcast
	mqttBroker := utils.MqttBroker
	return mqttBroker.Subscribe(ctx, topic, onMessageReceived)
}

func onMessageReceived(client mqtt.Client, msg mqtt.Message) {
	ctx := context.Background()
	payload := msg.Payload()
	payloadStr := string(payload)

	log.Info(ctx).Msg("Received message [" + payloadStr + "]")
}
