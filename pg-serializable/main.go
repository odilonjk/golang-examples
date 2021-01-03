package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// inicia conexao com banco postgres
	host := os.Getenv("DB_SERVER")
	connStr := "postgres://booking_app:pg@" + host + "/booking?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// mapeamento do endpoint
	http.HandleFunc("/bookings", func(w http.ResponseWriter, r *http.Request) {
		bookingHandler(w, r, db)
	})

	// expoe a aplicacao na porta 8080
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Erro ao iniciar:", err.Error())
	}
}

func bookingHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// obtem os parametros da requisicao
	// para simplificar o exemplo, nao existem validacoes alem da regra de negocio do README
	err := r.ParseForm()
	if err != nil {
		log.Fatal("Erro ao obter parâmetros:", err.Error())
	}
	startDate := r.Form.Get("start_date")
	endDate := r.Form.Get("end_date")

	// tenta criar reserva, eh aqui que fica o que interessa do exemplo
	id, err := createBooking(db, startDate, endDate)
	if err != nil {
		log.Println("Nao foi possivel criar reserva")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	log.Println("Reserva criada")
	fmt.Fprint(w, "Reserva realizada: ", id)
}
