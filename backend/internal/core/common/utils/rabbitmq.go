package utils

import (
	"context"
	"encoding/json"
	"fernandoglatz/home-management/internal/core/common/utils/constants"
	"fernandoglatz/home-management/internal/core/common/utils/log"
	"fernandoglatz/home-management/internal/infrastructure/config"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const DIRECT_EXCHANGE = "direct"

var RabbitMqBroker RabbitMqBrokerType

type RabbitMqBrokerType struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

type RabbitMqMessageHandler func(queue string, delivery amqp.Delivery) error

func ConnectToRabbitMQ(ctx context.Context) error {
	log.Info(ctx).Msg("Connecting to RabbitMQ")

	uri := config.ApplicationConfig.Broker.RabbitMQ.Uri

	connection, err := amqp.Dial(uri)
	if err != nil {
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	RabbitMqBroker = RabbitMqBrokerType{
		Connection: connection,
		Channel:    channel,
	}

	log.Info(ctx).Msg("Connected to RabbitMQ!")

	return nil
}

func (rabbitMqBroker *RabbitMqBrokerType) Bind(ctx context.Context, queue string, routingKey string, exchange string) error {
	log.Info(ctx).Msg("Binding queue [" + queue + "] to exchange [" + exchange + "] with routing key [" + routingKey + "]")
	return rabbitMqBroker.Channel.QueueBind(queue, routingKey, exchange, false, nil)
}

func (rabbitMqBroker *RabbitMqBrokerType) CreateExchange(ctx context.Context, exchange string) error {
	log.Info(ctx).Msg("Created exchange [" + exchange + "]")
	return rabbitMqBroker.Channel.ExchangeDeclare(exchange, DIRECT_EXCHANGE, true, false, false, false, nil)
}

func (rabbitMqBroker *RabbitMqBrokerType) CreateQueue(ctx context.Context, queue string) error {
	log.Info(ctx).Msg("Created queue [" + queue + "]")
	_, err := rabbitMqBroker.Channel.QueueDeclare(queue, true, false, false, false, nil)

	return err
}

func (rabbitMqBroker *RabbitMqBrokerType) PublishExchange(ctx context.Context, exchange string, object any) error {
	return rabbitMqBroker.publish(ctx, constants.EMPTY, exchange, object)
}

func (rabbitMqBroker *RabbitMqBrokerType) PublishQueue(ctx context.Context, queue string, object any) error {
	return rabbitMqBroker.publish(ctx, queue, constants.EMPTY, object)
}

func (rabbitMqBroker *RabbitMqBrokerType) publish(ctx context.Context, queue string, exchange string, object any) error {
	jsonData, err := json.Marshal(object)
	if err != nil {
		return err
	}

	publishing := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonData,
		Timestamp:   time.Now(),
	}

	json := string(jsonData)
	if queue == constants.EMPTY {
		log.Info(ctx).Msg("Published message [" + json + "] to queue [" + queue + "]")
	} else {
		log.Info(ctx).Msg("Published message [" + json + "] to exchange [" + exchange + "]")
	}

	return rabbitMqBroker.Channel.PublishWithContext(ctx, exchange, queue, false, false, publishing)
}

func (rabbitMqBroker *RabbitMqBrokerType) Subscribe(ctx context.Context, queue string, callback RabbitMqMessageHandler) error {
	log.Info(ctx).Msg("Subscribing to queue [" + queue + "]")

	messages, err := rabbitMqBroker.Channel.Consume(queue, constants.EMPTY, false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range messages {
			err = callback(queue, delivery)
			if err != nil {
				delivery.Nack(false, true)
			} else {
				delivery.Ack(false)
			}
		}
	}()

	log.Info(ctx).Msg("Subscribed to queue [" + queue + "]!")

	return nil
}
