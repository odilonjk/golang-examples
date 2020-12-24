package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	out1 := make(chan string)
	out2 := make(chan string)
	var wg sync.WaitGroup

	// the WaitGroup must know how much processes it needs to wait to finish
	// change to 1 and one of the processes below will stop earlier than expected
	wg.Add(2)

	// using channels to read the result of processes
	go process("order", out1)
	go process("refund", out2)

	go func() {
		for {
			msg, open := <-out1
			// if the channel is closed, it will leave the loop
			if !open {
				break
			}
			fmt.Println(msg)
		}
		wg.Done()
	}()

	// the exact thing than above, however with less code
	go func() {
		for msg := range out2 {
			fmt.Println(msg)
		}
		wg.Done()
	}()

	// wait until all processes are finished
	wg.Wait()
}

func process(item string, out chan string) {
	// always close the channel on the producer side, otherwise it might Panic!
	defer close(out)

	for i := 1; i <= 5; i++ {
		if item == "refund" {
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second / 2)
		out <- item
	}
}
