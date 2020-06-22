package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/google/uuid"
	"github.com/odilonjk/golang-examples/processamento-assincrono/estoque/pkg/rabbit"
	"github.com/streadway/amqp"
)

const createQueue = "create.productionorder"

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

	ok := true
	for ok {
		msg, hasMore := rabbit.Get(ch, createQueue)
		ok = hasMore

		err := processMsg(msg)
		if err != nil {
			handleRetry(ch, msg)
		}
		msg.Ack(false)
	}

	log.Println("All production orders where processed")
}

// processMsg will send to retry queue production orders with quantity over 500000
func processMsg(msg amqp.Delivery) (err error) {
	if msg.Body != nil {
		var m message
		err := json.Unmarshal(msg.Body, &m)
		failOnError(err, "Failed to decode message")
		if m.Quantity > 500000 {
			err = errors.New("Maximum quantity exceeded")
		}
		log.Println("Created production order for ref. code: " + m.RefCode.String())
	}
	return
}

// handleRetry will calculate delay, retry attempts, and finally will publish it
func handleRetry(ch *amqp.Channel, m amqp.Delivery) {
	delay := calculateDelay(m.Headers)
	log.Println(fmt.Sprintf("Delayed for %dms", delay))

	retryCount := calculateRetry(m.Headers)
	log.Println(fmt.Sprintf("Retrying for %d time", retryCount))

	// required header with delay and retry counter
	headers := make(amqp.Table)
	headers["x-delay"] = delay
	headers["x-retry-count"] = retryCount

	rabbit.Publish(ch, m, headers)
}

// calculateDelay adds 5000ms to last delay found on message header
func calculateDelay(h amqp.Table) int32 {
	var d int32
	if h["x-delay"] != nil {
		n := h["x-delay"].(int32)
		d = int32(math.Abs(float64(n)))
	}
	return d + 5000
}

// calculateRetry adds 1 to last retry counter found on message header
func calculateRetry(h amqp.Table) int32 {
	var r int32
	lastCount := h["x-retry-count"]
	if lastCount != nil {
		r = lastCount.(int32)
	}
	return r + 1
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
