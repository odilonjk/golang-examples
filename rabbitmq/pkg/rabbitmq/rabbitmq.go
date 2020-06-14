package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

// NewDefaultChannel creates a connection and returns a valid channel to RabbitMQ
func NewDefaultChannel() amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open channel")
	}
	return *ch
}

// NewQueue declares and returns a RabbitMQ queue
func NewQueue(ch *amqp.Channel, n string) amqp.Queue {
	q, err := ch.QueueDeclare(
		n,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
	return q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
