package game

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/odilonjk/pedra-papel-tesoura-llm/internal/llm"
)

type Game struct {
	model llm.LLM
	// player1score int
	// player2score int
}

func New(model llm.LLM) *Game {
	return &Game{model: model}
}

func (g *Game) Start(ctx context.Context) {
	fmt.Println("Iniciando novo jogo!")

	for {
		stop := nextRun(ctx, g.model)
		if stop {
			break
		}
		fmt.Printf("\nProxima rodada!")

	}
	fmt.Printf("\nJogo finalizado.")
}

func nextRun(ctx context.Context, model llm.LLM) (stop bool) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	for i := range 2 {
		player := i + 1
		go func(ctx context.Context, wg *sync.WaitGroup, model llm.LLM, player int) {
			defer wg.Done()
			makeMove(ctx, model, player)
		}(ctx, wg, model, player)
	}
	wg.Wait()
	for {
		fmt.Printf("\nContinue? [y/n]\n")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		v := strings.ToLower(input.Text())
		switch v {
		case "y":
			stop = false
			return
		case "n":
			stop = true
			return
		default:
			fmt.Printf("\nPlease, select a valid input.")
		}
	}
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
