package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/odilonjk/golang-examples/processamento-assincrono/estoque/pkg/rabbit"
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
}

func main() {
	log.Println("Preparing to sent message to queue")

	ch := rabbit.NewDefaultChannel()
	defer ch.Close()

	// required args to use Delayed Message plugin
	exchangeArgs := make(amqp.Table)
	exchangeArgs["x-delayed-type"] = "direct"

	rabbit.DeclareExchange(ch, productionOrderExchange, "x-delayed-message", exchangeArgs)
	rabbit.DeclareExchange(ch, deadLetterExchange, "direct", nil)

	// config args for queue
	queueArgs := make(amqp.Table)
	queueArgs["x-max-length"] = 3
	queueArgs["x-dead-letter-exchange"] = deadLetterExchange

	rabbit.DeclareQueue(ch, createQueue, queueArgs)
	rabbit.DeclareQueue(ch, deadLetterQueue, nil)

	rabbit.BindQueue(ch, productionOrderExchange, createQueue, routingKey)
	rabbit.BindQueue(ch, deadLetterExchange, deadLetterQueue, routingKey)

	m := createMsgAsJSON()
	rabbit.Publish(ch, m, productionOrderExchange, routingKey)

	log.Println("Message sent to queue")
}

// createMsgAsJSON creates a message with random values and parses it to JSON
func createMsgAsJSON() (j []byte) {
	rand.Seed(time.Now().UnixNano())
	qty := rand.Intn(9999999)
	refCode := uuid.New()
	colorCode := uuid.New()
	msg := message{refCode, colorCode, qty}

	j, err := json.Marshal(msg)
	failOnError(err, "Failed to parse message into JSON")
	return
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
