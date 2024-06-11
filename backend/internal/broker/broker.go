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
const SUFFIX_EX = "-ex"
const SUFFIX_DLQ = "-dlq"
const SUFFIX_DLEX = "-dlex"
const CONFIG_QUEUE_EVENTS = "events"
const X_DEAD_LETTER_EXCHANGE = "x-dead-letter-exchange"

func Setup(ctx context.Context) error {
	err := utils.ConnectToMQTT(ctx)
	if err != nil {
		return err
	}

	err = utils.ConnectToRabbitMQ(ctx, onRabbitMqConnected)
	if err != nil {
		return err
	}

	mqttBroker := utils.MqttBroker
	topics := config.ApplicationConfig.Broker.Mqtt.Topics
	topicBroadcast := topics.Broadcast

	err = mqttBroker.Subscribe(ctx, topicBroadcast, onMqttBroadCastMessageReceived)
	if err != nil {
		return err
	}

	return nil
}

func onRabbitMqConnected(ctx context.Context) error {
	rabbitMqBroker := utils.RabbitMqBroker

	topics := config.ApplicationConfig.Broker.Mqtt.Topics
	queues := config.ApplicationConfig.Broker.RabbitMQ.Queues

	topicEvents := topics.Events
	queueEvents := queues[CONFIG_QUEUE_EVENTS]
	routingKeyEvents := getRoutingKey(topicEvents)

	queue := queueEvents.Name
	dlqQueue := queue + SUFFIX_DLQ
	dlExchange := queue + SUFFIX_DLEX
	requeueExchange := queueEvents.RequeueDelayExchange

	err := rabbitMqBroker.CreateQueue(ctx, dlqQueue, nil)
	if err != nil {
		return err
	}

	err = rabbitMqBroker.CreateExchange(ctx, dlExchange)
	if err != nil {
		return err
	}

	err = rabbitMqBroker.Bind(ctx, dlqQueue, constants.HASH, dlExchange)
	if err != nil {
		return err
	}

	queueArgs := make(map[string]any)
	queueArgs[X_DEAD_LETTER_EXCHANGE] = dlExchange
	err = rabbitMqBroker.CreateQueue(ctx, queue, queueArgs)
	if err != nil {
		return err
	}

	err = rabbitMqBroker.CreateDelayedExchange(ctx, requeueExchange)
	if err != nil {
		return err
	}

	err = rabbitMqBroker.Bind(ctx, queue, queue, requeueExchange)
	if err != nil {
		return err
	}

	err = rabbitMqBroker.Bind(ctx, queue, routingKeyEvents, AMQ_TOPIC)
	if err != nil {
		return err
	}

	err = rabbitMqBroker.Subscribe(ctx, queueEvents, onRabbitMqEventMessageReceived)
	if err != nil {
		return err
	}

	return nil
}

func onRabbitMqEventMessageReceived(ctx context.Context, queue string, delivery amqp.Delivery) error {
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
