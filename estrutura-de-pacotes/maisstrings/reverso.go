// Pacote com funções de manipulação de strings que utilizam encode UTF-8.
// Utilizado para disponibilizar funções além das disponíveis no pacote padrão "strings".
package maisstrings

func Reverter(frase string) string {
	r := []rune(frase)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
