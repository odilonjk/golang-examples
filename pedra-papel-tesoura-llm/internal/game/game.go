package game

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/odilonjk/pedra-papel-tesoura-llm/internal/llm"
)

type Game struct {
	llm llm.LLM
	// player1score int
	// player2score int
}

func New(llm llm.LLM) *Game {
	return &Game{llm: llm}
}

func (g *Game) Start() {
	fmt.Println("Iniciando novo jogo!")

	wg := &sync.WaitGroup{}
	wg.Add(2)
	ctx := context.Background()

	for i := 0; i < 2; i++ {
		go func(ctx context.Context, wg *sync.WaitGroup, llm llm.LLM, player int) {
			defer wg.Done()
			makeMove(ctx, llm, player)
		}(ctx, wg, g.llm, i+1)
	}
	wg.Wait()
}

func makeMove(ctx context.Context, l llm.LLM, player int) {
	move, err := l.GetNextMove(ctx)
	if err != nil {
		fmt.Printf("\nJogador %d falhou ao fazer jogada: %s", player, err.Error())
		return
	}
	normalizedMove := normalize(move)
	if normalizedMove == "" {
		fmt.Printf("\nJogador %d nao fez jogada valida.", player)
	}
	fmt.Printf("\nJogador %d selecionou: %s", player, normalizedMove)
}

// I need to normalize since the model I'm running locally
// hallucinates A LOT.
func normalize(s string) string {
	ls := strings.ToLower(s)
	switch {
	case strings.Contains(ls, "pedra"):
		return "pedra"
	case strings.Contains(ls, "papel"):
		return "papel"
	case strings.Contains(ls, "tesoura"):
		return "tesoura"
	default:
		return ""
	}
}
