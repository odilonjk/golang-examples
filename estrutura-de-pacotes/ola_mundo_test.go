package main

import "testing"

const olaMundo = "Olá, mundo!"
const olaRevertido = "!odnum ,álO"

func TestOlaMundo(t *testing.T) {
	esperado := olaMundo
	recebeu := OlaMundo()
	if esperado != recebeu {
		t.Errorf("esperava %s mas recebeu %s", esperado, recebeu)
	}
}

func TestOlaMundoRevertido(t *testing.T) {
	esperado := olaRevertido
	recebeu := OlaMundoRevertido()
	if esperado != recebeu {
		t.Errorf("esperava %s mas recebeu %s", esperado, recebeu)
	}
}
