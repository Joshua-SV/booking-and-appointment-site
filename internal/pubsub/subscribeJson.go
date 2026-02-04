package pubsub

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	return Subscribe[T](conn, exchange, queueName, key, queueType, handler,
		func(data_bytes []byte) (T, error) {
			var val T
			// unmarshal the JSON message into the specified type
			err := json.Unmarshal(data_bytes, &val)
			if err != nil {
				return val, err
			}
			return val, nil
		})
}
