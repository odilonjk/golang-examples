package main

import (
	"log"

	nats "github.com/nats-io/nats.go"
)

func main() {
	opts := []nats.Option{nats.Name("NATS Sample Publisher")}
	log.Printf("NATS port %s", nats.DefaultURL)

	nc, err := nats.Connect(nats.DefaultURL, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	subj := "first.subject"
	msg := "Hello World"
	nc.Publish(subj, []byte(msg))
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subj, msg)
	}

}
