package util

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	rabbitmq "github.com/wagslane/go-rabbitmq"
)

func Publish(amqp_url string, exchange string, routing_key string, message WireguardMessage) error {
	publisher, err := rabbitmq.NewPublisher(
		amqp_url, amqp.Config{},
		rabbitmq.WithPublisherOptionsLogging,
	)
	if err != nil {
		return err
	}
	data, _ := json.Marshal(message)
	err = publisher.Publish(
		data,
		[]string{routing_key},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange(exchange),
	)
	if err != nil {
		return err
	}
	return nil
}
