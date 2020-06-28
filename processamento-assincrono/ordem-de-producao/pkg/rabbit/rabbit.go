package rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

// NewDefaultChannel creates a connection and returns a valid channel to RabbitMQ
func NewDefaultChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open channel")
	}
	return ch
}

// Publish sends the message to the given exchange and routing key
func Publish(ch *amqp.Channel, e, r string, b []byte, h amqp.Table) {
	err := ch.Publish(
		e,
		r,
		false,
		false,
		amqp.Publishing{
			Headers:     h,
			ContentType: "application/json",
			Body:        b,
		})
	failOnError(err, "Failed to publish on exchange "+e)
}

// Consume returns a consumer for the given queue
func Consume(ch *amqp.Channel, q string) (d <-chan amqp.Delivery) {
	d, err := ch.Consume(
		q,
		"production-order",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to consume from queue "+q)
	return
}

// failOnError log and exit in case of error
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
