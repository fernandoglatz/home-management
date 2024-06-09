package utils

import (
	"context"
	"fernandoglatz/home-management/internal/core/common/utils/log"
	"fernandoglatz/home-management/internal/infrastructure/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMqBroker RabbitMqBrokerType

type RabbitMqBrokerType struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

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
