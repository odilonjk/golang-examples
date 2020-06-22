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

// DeclareExchange declares an exchange with given name, type and args
func DeclareExchange(ch *amqp.Channel, n, t string, args amqp.Table) {
	err := ch.ExchangeDeclare(
		n,
		t,
		true,
		false,
		false,
		false,
		args,
	)
	failOnError(err, "Failed to declare exchange "+n)
}

// DeclareQueue declares a queue with given name and args
func DeclareQueue(ch *amqp.Channel, n string, args amqp.Table) {
	_, err := ch.QueueDeclare(
		n,
		true, false,
		false,
		false,
		args,
	)
	failOnError(err, "Failed to declare queue "+n)
}

// BindQueue creates the binding between a queue and an exchange given a routing key
func BindQueue(ch *amqp.Channel, e, q, k string) {
	err := ch.QueueBind(
		q,
		k,
		e,
		false,
		nil)
	failOnError(err, "Failed to bind "+e+" with "+q)
}

// Publish sends the JSON message to the given exchange and routing key
func Publish(ch *amqp.Channel, m []byte, e, k string) {
	err := ch.Publish(
		e,
		k,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        m,
		})
	failOnError(err, "Failed to publish on exchange "+e)
}

// failOnError log and exit in case of error
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
