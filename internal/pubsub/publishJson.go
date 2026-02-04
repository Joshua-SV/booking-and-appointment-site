package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishJSON publishes a JSON-encoded message to the specified exchange with the given routing key.
func PublishJSON[T any](conn *amqp.Connection, exchange, key string, val T) error {
	// create a channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Marshal the value to JSON bytes
	body, err := json.Marshal(val)
	if err != nil {
		return err
	}

	// Publish the message to the exchange with the routing key
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return err
	}

	return nil
}
