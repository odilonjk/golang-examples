package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", imprimirOlaMundo)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Utilizando porta padrão %s", port)
	}

	log.Printf("Serviço disponível na porta %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}

func imprimirOlaMundo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Olá mundo! Eu consegui!!! :D")
}
