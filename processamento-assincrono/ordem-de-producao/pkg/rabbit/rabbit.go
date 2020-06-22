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
func Publish(ch *amqp.Channel, d amqp.Delivery, h amqp.Table) {
	err := ch.Publish(
		d.Exchange,
		d.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Headers:     h,
			ContentType: "application/json",
			Body:        d.Body,
		})
	failOnError(err, "Failed to publish on exchange "+d.Exchange)
}

// Get the delivery from queue
func Get(ch *amqp.Channel, q string) (amqp.Delivery, bool) {
	msg, hasMore, err := ch.Get(
		q,
		false,
	)
	failOnError(err, "Failed to get message from queue "+q)
	return msg, hasMore
}

// failOnError log and exit in case of error
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
