package pubsub

import amqp "github.com/rabbitmq/amqp091-go"

type SimpleQueueType int

const (
	TransientQueue SimpleQueueType = iota // start at 0 this queue will not survive server restarts
	DurableQueue                          // this queue will survive server restarts
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	// create a channel
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	// create a queue using the channel
	queue, err := channel.QueueDeclare(queueName, queueType == DurableQueue, queueType == TransientQueue, queueType == TransientQueue, false, amqp.Table{
		"x-dead-letter-exchange": "appointment_dlx", // sends messages to dead-letter exchange upon rejection
	})
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	// bind the queue to the exchange with the routing key
	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return channel, queue, nil
}
