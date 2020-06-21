package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

const (
	productionOrderExchange = "productionorder"
	deadLetterExchange      = "productionorder.dlx"
	createQueue             = "create.productionorder"
	deadLetterQueue         = "create.productionorder.dlx"
	routingKey              = "create"
)

// message sent to RabbitMQ
type message struct {
	RefCode   uuid.UUID
	ColorCode uuid.UUID
	Quantity  int
	EventType string
}

func main() {
	log.Println("Preparing to sent message to queue")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	exchangeArgs := make(amqp.Table)
	exchangeArgs["x-delayed-type"] = "direct"
	err = declareExchage(ch, productionOrderExchange, "x-delayed-message", exchangeArgs)
	failOnError(err, "Failed to declare exchange "+productionOrderExchange)

	err = declareExchage(ch, deadLetterExchange, "direct", nil)
	failOnError(err, "Failed to declare exchange "+deadLetterExchange)

	queueArgs := make(amqp.Table)
	queueArgs["x-max-length"] = 3
	queueArgs["x-dead-letter-exchange"] = deadLetterExchange

	err = declareQueue(ch, createQueue, queueArgs)
	failOnError(err, "Failed to declare queue")
	err = declareQueue(ch, deadLetterQueue, nil)
	failOnError(err, "Failed to declare queue")

	err = bindQueue(ch, productionOrderExchange, createQueue)
	failOnError(
		err,
		fmt.Sprintf("Failed to bind exchange %s with queue %s", productionOrderExchange, createQueue),
	)
	err = bindQueue(ch, deadLetterExchange, deadLetterQueue)
	failOnError(
		err,
		fmt.Sprintf("Failed to bind exchange %s with queue %s", productionOrderExchange, createQueue),
	)

	m := createMsgAsJSON()
	err = publishMessage(ch, m)
	failOnError(err, "Failed to publish message")

	log.Println("Message sent to queue")
}

// declareExchange creates a exchange
func declareExchage(ch *amqp.Channel, n, t string, args amqp.Table) error {
	return ch.ExchangeDeclare(
		n,
		t,
		true,
		false,
		false,
		false,
		args,
	)
}

// newQueue declares a new queue
func declareQueue(ch *amqp.Channel, n string, args amqp.Table) (err error) {
	_, err = ch.QueueDeclare(
		n,
		true, false,
		false,
		false,
		args,
	)
	return
}

// bindQueue creates the binding between a queue and an exchange
func bindQueue(ch *amqp.Channel, exchange, queue string) error {
	return ch.QueueBind(
		queue,
		routingKey,
		exchange,
		false,
		nil)
}

// publishMessage sends the JSON message to the given exchange
func publishMessage(ch *amqp.Channel, m []byte) error {
	return ch.Publish(
		productionOrderExchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        m,
		})
}

// createMsgAsJSON creates a message with random values and parses it to JSON
func createMsgAsJSON() (j []byte) {
	rand.Seed(time.Now().UnixNano())
	qty := rand.Intn(9999999)
	refCode := uuid.New()
	color := uuid.New()
	msg := message{refCode, color, qty, "CREATE"}

	j, err := json.Marshal(msg)
	failOnError(err, "Failed to parse message into JSON")
	return
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
