package main

import (
	"fmt"

	"github.com/odilonjk/golang-examples/estrutura-de-pacotes/maisstrings"
)

func main() {
	fmt.Println(OlaMundo())
	fmt.Println(OlaMundoRevertido())
}

func OlaMundo() string {
	return "Olá, mundo!"
}

func OlaMundoRevertido() string {
	return maisstrings.Reverter(OlaMundo())
}
