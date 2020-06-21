package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"

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
	log.Println("Production Order initialized")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	ok := true
	for ok {
		msg, hasMore, err := ch.Get(
			createQueue,
			false,
		)
		ok = hasMore
		failOnError(err, "Failed to get message")

		if msg.Body != nil {
			var m message
			err = json.Unmarshal(msg.Body, &m)
			failOnError(err, "Failed to decode message")
			if m.Quantity > 500000 {
				log.Println("Error processing ref. code: " + m.RefCode.String())
				err := publishRetry(ch, msg)
				failOnError(err, fmt.Sprintf("Failed to retry %v", m))
			}
			log.Println("Created production order for ref. code: " + m.RefCode.String())
			msg.Ack(false)
		}
	}

	log.Println("All production orders where processed")
}

func publishRetry(ch *amqp.Channel, d amqp.Delivery) error {

	var delay int32
	if d.Headers["x-delay"] != nil {
		n := d.Headers["x-delay"].(int32)
		delay = int32(math.Abs(float64(n)))
	}
	delay += 5000
	log.Println(fmt.Sprintf("Delayed for %dms", delay))

	var retry int32
	if d.Headers["x-retry-count"] != nil {
		retry = d.Headers["x-retry-count"].(int32)
	}
	retry++
	log.Println(fmt.Sprintf("Retrying for %d time", retry))

	headers := make(amqp.Table)
	headers["x-delay"] = delay
	headers["x-retry-count"] = retry

	return ch.Publish(
		d.Exchange,
		d.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Headers:     headers,
			ContentType: "application/json",
			Body:        d.Body,
		})
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
