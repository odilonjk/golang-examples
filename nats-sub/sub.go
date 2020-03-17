package main

import (
	"flag"
	"log"
	"runtime"

	nats "github.com/nats-io/nats.go"
)

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))
}

func main() {

	var showTime = flag.Bool("t", false, "Display timestamps")

	opts := []nats.Option{nats.Name("NATS Sample Subscriber")}
	log.Printf("NATS port %s", nats.DefaultURL)

	nc, err := nats.Connect(nats.DefaultURL, opts...)
	if err != nil {
		log.Fatal(err)
	}

	subj := "first.subject"
	i := 0
	nc.Subscribe(subj, func(msg *nats.Msg) {
		i++
		printMsg(msg, i)
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]", subj)
	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	runtime.Goexit()

}
