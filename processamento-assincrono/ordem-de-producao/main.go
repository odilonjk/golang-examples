package main

import (
	"encoding/json"
	"errors"
	"log"
	"math"

	"github.com/google/uuid"
	"github.com/odilonjk/golang-examples/processamento-assincrono/ordem-de-producao/pkg/rabbit"
	"github.com/streadway/amqp"
)

const (
	productionOrderExchange = "productionorder"
	deadLetterExchange      = "productionorder.dlx"
	createQueue             = "create.productionorder"
	routingKey              = "create"
)

// message sent to RabbitMQ
type message struct {
	RefCode   uuid.UUID
	ColorCode uuid.UUID
	Quantity  int
}

func main() {
	log.Println("Production Order initialized")

	ch := rabbit.NewDefaultChannel()
	defer ch.Close()

	msgs := rabbit.Consume(ch, createQueue)

	forever := make(chan bool)
	go func() {
		for m := range msgs {
			err := processMsg(m)
			if err != nil {
				handleFailedMsg(ch, m)
			}
			m.Ack(false)
		}
	}()

	log.Println("Waiting messages to be delivered")
	<-forever
}

// processMsg will send to retry queue production orders with quantity over 500000
func processMsg(msg amqp.Delivery) (err error) {
	if msg.Body != nil {
		var m message
		err = json.Unmarshal(msg.Body, &m)
		failOnError(err, "Failed to decode message")
		if m.Quantity > 500 {
			err = errors.New("Maximum quantity exceeded")
		}
		log.Println("Created production order for ref. code: " + m.RefCode.String())
	}
	return
}

// handleFailedMsg will calculate delay, retry attempts, and finally will publish it
func handleFailedMsg(ch *amqp.Channel, m amqp.Delivery) {
	retryCount := getRetry(m.Headers)
	headers := make(amqp.Table)
	if retryCount > 2 {
		rabbit.Publish(ch, deadLetterExchange, routingKey, m.Body, nil)
	} else {
		headers["x-delay"] = getDelay(m.Headers) + 5000
		headers["x-retry-count"] = retryCount + 1
		rabbit.Publish(ch, productionOrderExchange, routingKey, m.Body, headers)
	}
}

// getDelay returns the delay value from the header
func getDelay(h amqp.Table) int32 {
	var d int32
	if h["x-delay"] != nil {
		n := h["x-delay"].(int32)
		d = int32(math.Abs(float64(n)))
	}
	return d
}

// getRetry returns the retry count from the header
func getRetry(h amqp.Table) int32 {
	var r int32
	lastCount := h["x-retry-count"]
	if lastCount != nil {
		r = lastCount.(int32)
	}
	return r
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
