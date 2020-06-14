package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/odilonjk/golang-examples/rabbitmq/pkg/rabbitmq"
	"github.com/streadway/amqp"
)

// Message sent to RabbitMQ, containing reference code and event type
type Message struct {
	RefCode   uuid.UUID
	EventType string
}

func main() {
	log.Println("Publisher initialized")

	ch := rabbitmq.NewDefaultChannel()
	defer ch.Close()
	q := rabbitmq.NewQueue(&ch, "fallback")

	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		id := uuid.New()
		msg := Message{id, "DELETE"}
		json, err := json.Marshal(msg)
		failOnError(err, "Failed to parse message into JSON")
		wg.Add(1)
		go func(wg *sync.WaitGroup, json []byte) {
			defer wg.Done()
			err = ch.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        json,
				})
			failOnError(err, "Failed to publish message")
		}(&wg, json)
	}
	wg.Wait()

	log.Println("Publisher shutting down")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
