package main

import "testing"

const olaMundo = "Olá, mundo!"

func TestOlaMundo(t *testing.T) {
	espera := olaMundo
	recebeu := OlaMundo()
	if espera != recebeu {
		t.Errorf("esperava %s mas recebeu %s", espera, recebeu)
	}
}
