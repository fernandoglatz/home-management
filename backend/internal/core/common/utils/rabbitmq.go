package utils

import (
	"context"
	"encoding/json"
	"fernandoglatz/home-management/internal/core/common/utils/constants"
	"fernandoglatz/home-management/internal/core/common/utils/log"
	"fernandoglatz/home-management/internal/infrastructure/config"
	"fmt"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

const HEADER_RECEIVED_COUNT = "received-count"
const X_DELAYED_MESSAGE = "x-delayed-message"
const X_DELAYED_TYPE = "x-delayed-type"
const X_DELAY = "x-delay"
const DIRECT = "direct"

var RabbitMqBroker RabbitMqBrokerType

type RabbitMqBrokerType struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

type RabbitMqMessageHandler func(ctx context.Context, queue string, delivery amqp.Delivery) error

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
	log.Info(ctx).Msg("Creating exchange [" + exchange + "]")
	return rabbitMqBroker.Channel.ExchangeDeclare(exchange, amqp.ExchangeFanout, true, false, false, false, nil)
}

func (rabbitMqBroker *RabbitMqBrokerType) CreateDelayedExchange(ctx context.Context, exchange string) error {
	log.Info(ctx).Msg("Creating delayed exchange [" + exchange + "]")

	args := make(map[string]any)
	args[X_DELAYED_TYPE] = DIRECT
	return rabbitMqBroker.Channel.ExchangeDeclare(exchange, X_DELAYED_MESSAGE, true, false, false, false, args)
}

func (rabbitMqBroker *RabbitMqBrokerType) CreateQueue(ctx context.Context, queue string, args map[string]any) error {
	log.Info(ctx).Msg("Creating queue [" + queue + "]")

	if args == nil {
		args = make(map[string]any)
	}

	args[amqp.QueueTypeArg] = amqp.QueueTypeQuorum
	_, err := rabbitMqBroker.Channel.QueueDeclare(queue, true, false, false, false, args)

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
	}

	json := string(jsonData)
	if queue == constants.EMPTY {
		log.Info(ctx).Msg("Publishing message [" + json + "] to queue [" + queue + "]")
	} else {
		log.Info(ctx).Msg("Publishing message [" + json + "] to exchange [" + exchange + "]")
	}

	return rabbitMqBroker.Channel.PublishWithContext(ctx, exchange, queue, false, false, publishing)
}

func (rabbitMqBroker *RabbitMqBrokerType) Subscribe(ctx context.Context, config config.Queue, callback RabbitMqMessageHandler) error {
	queue := config.Name
	log.Info(ctx).Msg("Subscribing to queue [" + queue + "]")

	messages, err := rabbitMqBroker.Channel.Consume(queue, constants.EMPTY, false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range messages {
			traceMap := make(map[string]any)
			traceMap[constants.MESSAGE_ID] = delivery.DeliveryTag

			callbackCtx := context.Background()
			callbackCtx = context.WithValue(callbackCtx, constants.TRACE_MAP, traceMap)

			err = callback(callbackCtx, queue, delivery)
			rabbitMqBroker.handleCallbackResult(callbackCtx, config, delivery, err)
		}
	}()

	log.Info(ctx).Msg("Subscribed to queue [" + queue + "]!")

	return nil
}

func (rabbitMqBroker *RabbitMqBrokerType) handleCallbackResult(ctx context.Context, config config.Queue, delivery amqp.Delivery, err error) {
	if err == nil {
		delivery.Ack(false)

	} else {
		queue := config.Name
		maximumReceives := config.MaximumReceives

		if delivery.Headers == nil {
			delivery.Headers = make(map[string]any)
		}

		count := constants.ZERO
		countObj := delivery.Headers[HEADER_RECEIVED_COUNT]
		if countObj != nil {
			count, _ = strconv.Atoi(countObj.(string))
		}
		count++

		if count >= maximumReceives {
			json := string(delivery.Body)
			message := fmt.Sprintf("Moving message [%s] to dlq for queue [%s] with count [%d]", json, queue, count)
			log.Info(ctx).Msg(message)

			err := delivery.Nack(false, false)
			if err != nil {
				message := fmt.Sprintf("Error on moving message [%s] to dlq for queue [%s] with count [%d]", json, queue, count)

				log.Error(ctx).Msg(message)
			}

		} else {
			delivery.Ack(false)
			delivery.Headers[HEADER_RECEIVED_COUNT] = strconv.Itoa(count)
			rabbitMqBroker.requeue(ctx, config, delivery)
		}
	}
}

func (rabbitMqBroker *RabbitMqBrokerType) requeue(ctx context.Context, config config.Queue, delivery amqp.Delivery) error {
	queue := config.Name
	exchange := config.RequeueDelayExchange

	publishing := amqp.Publishing{
		Headers:         delivery.Headers,
		ContentType:     delivery.ContentType,
		ContentEncoding: delivery.ContentEncoding,
		DeliveryMode:    delivery.DeliveryMode,
		Priority:        delivery.Priority,
		CorrelationId:   delivery.CorrelationId,
		ReplyTo:         delivery.ReplyTo,
		Expiration:      delivery.Expiration,
		MessageId:       delivery.MessageId,
		Timestamp:       delivery.Timestamp,
		Type:            delivery.Type,
		UserId:          delivery.UserId,
		AppId:           delivery.AppId,
		Body:            delivery.Body,
	}

	count := publishing.Headers[HEADER_RECEIVED_COUNT]
	publishing.Headers[X_DELAY] = config.RequeueDelay.Milliseconds()

	json := string(delivery.Body)
	message := fmt.Sprintf("Requeuing message [%s] to exchange [%s] with count [%s]", json, exchange, count)
	log.Info(ctx).Msg(message)

	return rabbitMqBroker.Channel.PublishWithContext(ctx, exchange, queue, false, false, publishing)
}
