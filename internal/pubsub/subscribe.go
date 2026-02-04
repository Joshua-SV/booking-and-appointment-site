package pubsub

import amqp "github.com/rabbitmq/amqp091-go"

type AckType int

const (
	Ack AckType = iota
	Nack_discard
	Nack_requeue
)

func Subscribe[T any](conn *amqp.Connection,
	exchangeName,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
) error {
	// declare and bind a new queue to the exchange
	channel, queue, err := DeclareAndBind(conn, exchangeName, queueName, key, queueType)
	if err != nil {
		return err
	}

	// modify the prefetch count to 10 messages per consumer of the same queue
	err = channel.Qos(10, 0, false)
	if err != nil {
		return err
	}

	// start consuming messages from the queue
	messages, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	// start a goroutine to handle incoming messages
	go func() {
		for msg := range messages {
			// unmarshal the message
			val, err := unmarshaller(msg.Body)
			if err != nil {
				// if unmarshaling fails, nack the message and continue
				msg.Nack(false, false)
				continue
			}

			// call the handler function with the unmarshaled value
			ackResp := handler(val)

			if ackResp == Ack {
				// acknowledge the message after successful handling
				msg.Ack(false)
			} else if ackResp == Nack_discard {
				// nack the message and discard it (send to dead-letter queue if configured)
				msg.Nack(false, false)
			} else if ackResp == Nack_requeue {
				// nack the message and requeue it (resend to the queue)
				msg.Nack(false, true)
			}
		}
	}()

	return nil
}
