package strutil

import (
	"testing"
)

const fraseRevertida = "!iugesnoc uE"

func TestReverter(t *testing.T) {
	esperado := fraseRevertida
	recebeu := Reverter("Eu consegui!")
	if esperado != recebeu {
		t.Errorf("esperava %s mas recebeu %s", esperado, recebeu)
	}
}
