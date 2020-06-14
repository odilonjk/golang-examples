package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/odilonjk/golang-examples/rabbitmq/pkg/rabbitmq"
)

// Message consumed from RabbitMQ queue
type Message struct {
	RefCode   uuid.UUID
	EventType string
}

func main() {
	log.Println("Subscriber initialized")
	start := time.Now()

	ch := rabbitmq.NewDefaultChannel()

	q := rabbitmq.NewQueue(&ch, "fallback")
	var count int64
	ok := true
	for ok {
		msg, hasMore, err := ch.Get(
			q.Name,
			false,
		)
		ok = hasMore
		failOnError(err, "Failed to register a consumer")

		if msg.Body != nil {
			var m Message
			err = json.Unmarshal(msg.Body, &m)
			failOnError(err, "Failed to decode message")
			log.Println(fmt.Sprintf("Received message: %v", m))
			count++
			msg.Ack(false)
		}

	}

	elapsed := time.Since(start)
	log.Printf("Consumed %d in %dms", count, elapsed.Milliseconds())
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
