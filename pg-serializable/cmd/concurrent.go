package main

import (
	"log"
	"net/http"
	"net/url"
	"sync"
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
	out := make(chan bool, 6)
	var wg sync.WaitGroup

	// realiza 6 chamadas simultaneas para as instancias da aplicacao
	// sem o tratamento correto, resultaria em reservas sobrepostas
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go call(ports[i], out)
	}

	func() {
		wg.Wait()
		close(out)
	}()
}

// call realiza a chamada HTTP para gerar uma reserva
func call(p string, out chan bool) {

	path := "http://localhost:" + p + "/bookings"
	_, err := http.PostForm(path, url.Values{"start_date": {"2020-01-07"}, "end_date": {"2020-01-10"}})
	if err != nil {
		log.Fatal("Erro ao realizar requisicao de criacao de reserva: ", err.Error())
	}
	log.Println("Requisicao realizada para porta", p)
	out <- true

}
