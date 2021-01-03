package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var ports = [...]string{
	"8080",
	"8081",
	"8082",
	"8083",
	"8084",
	"8085",
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	out1 := make(chan bool)
	out2 := make(chan bool)
	out3 := make(chan bool)
	out4 := make(chan bool)
	out5 := make(chan bool)
	out6 := make(chan bool)
	var wg sync.WaitGroup

	wg.Add(6)

	// realiza 6 chamadas simultaneas para as instancias da aplicacao
	go call(out1)
	go call(out2)
	go call(out3)
	go call(out4)
	go call(out5)
	go call(out6)

	go waitBooking(&wg, out1)
	go waitBooking(&wg, out2)
	go waitBooking(&wg, out3)
	go waitBooking(&wg, out4)
	go waitBooking(&wg, out5)
	go waitBooking(&wg, out6)

	wg.Wait()
}

func call(out chan bool) {
	defer close(out)

	// as portas para a requisicao sao selecionadas aleatoriamente
	// pode acontecer de chamar repetidas vezes a mesma instancia
	rand.Seed(time.Now().UTC().UnixNano())
	n := rand.Int() % len(ports)
	path := "http://localhost:" + ports[n] + "/bookings"

	// tenta gerar a reserva para a mesma data
	_, err := http.PostForm(path, url.Values{"start_date": {"2020-01-07"}, "end_date": {"2020-01-10"}})
	if err != nil {
		log.Fatal("Erro ao realizar requisicao de criacao de reserva: ", err.Error())
	}
	log.Println(ports[n])
	out <- true

}

func waitBooking(wg *sync.WaitGroup, out chan bool) {
	for {
		_, open := <-out
		if !open {
			break
		}
	}
	wg.Done()
}
