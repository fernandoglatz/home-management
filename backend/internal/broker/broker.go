package broker

import (
	"context"
	"fernandoglatz/home-management/internal/core/common/utils"
	"fernandoglatz/home-management/internal/core/common/utils/constants"
	"fernandoglatz/home-management/internal/core/common/utils/log"
	"fernandoglatz/home-management/internal/core/entity"
	"fernandoglatz/home-management/internal/core/service"
	"fernandoglatz/home-management/internal/infrastructure/config"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	amqp "github.com/rabbitmq/amqp091-go"
)

const AMQ_TOPIC = "amq.topic"

func Setup(ctx context.Context) error {
	mqttBroker := utils.MqttBroker
	rabbitMqBroker := utils.RabbitMqBroker

	topics := config.ApplicationConfig.Broker.Mqtt.Topics
	queues := config.ApplicationConfig.Broker.RabbitMQ.Queues

	topicBroadcast := topics.Broadcast
	topicEvents := topics.Events
	queueEvents := queues.Events
	routingKeyEvents := getRoutingKey(topicEvents)

	err := mqttBroker.Subscribe(ctx, topicBroadcast, onMqttBroadCastMessageReceived)
	if err != nil {
		return err
	}

	err = rabbitMqBroker.CreateQueue(ctx, queueEvents)
	if err != nil {
		return err
	}

	err = rabbitMqBroker.Bind(ctx, queueEvents, routingKeyEvents, AMQ_TOPIC)
	if err != nil {
		return err
	}

	err = rabbitMqBroker.Subscribe(ctx, queueEvents, onRabbitMqEventMessageReceived)
	if err != nil {
		return err
	}

	return nil
}

func onRabbitMqEventMessageReceived(queue string, delivery amqp.Delivery) error {
	ctx := context.Background()
	defer log.HandlePanic(ctx)

	eventService := service.GetEventService[*entity.Event]()
	body := delivery.Body
	errw := eventService.ProcessMessage(ctx, body)

	if errw != nil {
		json := string(body)
		log.Error(ctx).PutTraceMap("json", json).Msg("Error on processing event message: " + errw.GetMessage())
		return errw.Error
	}

	return nil
}

func onMqttBroadCastMessageReceived(client mqtt.Client, msg mqtt.Message) {
	ctx := context.Background()
	payload := msg.Payload()
	payloadStr := string(payload)

	log.Info(ctx).Msg("Received MQTT message [" + payloadStr + "]")
}

func getRoutingKey(topic string) string {
	routingKey := strings.Replace(topic, constants.SLASH, constants.DOT, constants.MINUS_ONE)
	return strings.Replace(routingKey, constants.PLUS, constants.ASTERISK, constants.MINUS_ONE)
}
