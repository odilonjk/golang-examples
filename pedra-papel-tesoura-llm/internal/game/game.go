package game

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/odilonjk/pedra-papel-tesoura-llm/internal/llm"
	"github.com/odilonjk/pedra-papel-tesoura-llm/internal/models"
)

type status struct {
	turn       int
	p1score    int
	p2score    int
	p1move     string
	p2move     string
	p1messages []models.Message
	p2messages []models.Message
}

type Game struct {
	model  llm.LLM
	status status
}

func New(model llm.LLM) *Game {
	return &Game{
		model: model,
		status: status{
			turn:       1,
			p1messages: buildBaseMessages(),
			p2messages: buildBaseMessages(),
		},
	}
}

func (g *Game) Start(ctx context.Context) {
	fmt.Println("Iniciando novo jogo!")

	for {
		wg := &sync.WaitGroup{}
		wg.Add(2)

		movesCh := make(chan map[int]string, 2)
		g.callPlayers(ctx, movesCh)

		g.waitPlayers(wg, movesCh)

		g.computeResults()
		g.updateMessages()
		fmt.Printf("\nResultado atual:\nPlayer 1: %d\nPlayer 2: %d\n", g.status.p1score, g.status.p2score)
		if !anotherRun() {
			break
		}
		fmt.Printf("\nProxima rodada!")
	}
	fmt.Printf("\nJogo finalizado.")
}

func (g *Game) computeResults() {
	status := g.status
	winner := 0
	fmt.Printf("\nJogador 1 escolheu: %s\nJogador 2 escolheu: %s", status.p1move, status.p2move)
	if status.p1move == "" || status.p2move == "" {
		fmt.Printf("\nJogada invalida! Um ou ambos jogadores nao apresentaram sua opcao corretamente.")
		return
	}

	if status.p1move == status.p2move {
		fmt.Printf("\nEmpate!")
		return
	}

	if status.p1move == "pedra" {
		if status.p2move == "tesoura" {
			g.status.p1score = status.p1score + 1
		} else {
			g.status.p2score = status.p2score + 1
		}
	}

	if status.p1move == "papel" {
		if status.p2move == "pedra" {
			g.status.p1score = status.p1score + 1

		} else {
			g.status.p2score = status.p2score + 1
		}
	}

	if status.p1move == "tesoura" {
		if status.p2move == "papel" {
			g.status.p1score = status.p1score + 1

		} else {
			g.status.p2score = status.p2score + 1
			winner = 2
		}
	}
	fmt.Printf("\nO vencedor da rodada %d foi: Jogador %d", g.status.turn, winner)
	fmt.Printf("\nPlacar:\nJogador 1: %d\nJogador 2: %d", g.status.p1score, g.status.p2score)
}

func (g *Game) updateMessages() {
	g.status.p1messages = append(g.status.p1messages, buildUpdateMessages(g.status.p1move, g.status.p2move, g.status.p1score, g.status.p2score)...)
	g.status.p2messages = append(g.status.p2messages, buildUpdateMessages(g.status.p2move, g.status.p1move, g.status.p2score, g.status.p1score)...)
}

func (g *Game) waitPlayers(wg *sync.WaitGroup, movesCh chan map[int]string) {
	go func(wg *sync.WaitGroup, movesCh chan map[int]string) {
		for range 2 {
			for k, v := range <-movesCh {

				if k == 1 {
					g.status.p1move = v
				} else {
					g.status.p2move = v
				}
				wg.Done()
			}
		}
	}(wg, movesCh)
	wg.Wait()
}

func (g *Game) callPlayers(ctx context.Context, movesCh chan map[int]string) {
	for i := range 2 {
		player := i + 1
		var msgs []models.Message
		if player == 1 {
			msgs = g.status.p1messages
		} else {
			msgs = g.status.p2messages
		}
		go func(ctx context.Context, movesCh chan map[int]string, model llm.LLM, player int, msgs []models.Message) {
			fmt.Printf("\nJogador %d pensando...", player)
			makeMove(ctx, movesCh, model, player, msgs)
		}(ctx, movesCh, g.model, player, msgs)
	}
}

func makeMove(ctx context.Context, movesCh chan map[int]string, l llm.LLM, player int, msgs []models.Message) {
	move, err := l.GetNextMove(ctx, msgs)
	if err != nil {
		fmt.Printf("\nJogador %d falhou ao fazer jogada: %s", player, err.Error())
		movesCh <- map[int]string{player: ""}
	}
	normalizedMove := normalize(move)

	movesCh <- map[int]string{player: normalizedMove}
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

func anotherRun() (anotherRun bool) {
	for {
		fmt.Printf("\nContinuar? [s/n]\n")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		v := strings.ToLower(input.Text())
		switch v {
		case "s":
			anotherRun = true
			return
		case "n":
			anotherRun = false
			return
		default:
			fmt.Printf("\nPor favor, selecione um valor valido.")
		}
	}
}

func buildBaseMessages() []models.Message {
	return []models.Message{
		{Role: "system", Content: "You're playing the game 'Pedra, papel ou tesoura.'. You just need to select one of these values: pedra, or papel, or tesoura. Pedra is stronger than tesoura. Teroura is stronger than papel. Papel is stronger than pedra. You don't know which one your opponent will select. Try your best to win the round!"},
		{Role: "user", Content: "Your opponent is selecting one of the three options. Which one you select?"},
	}
}

func buildUpdateMessages(yourMove, opponentMove string, yourScore, otherScore int) []models.Message {
	if yourMove == "" {
		yourMove = "none"
	}
	if opponentMove == "" {
		opponentMove = "none"
	}
	return []models.Message{
		{
			Role:    "assistant",
			Content: yourMove,
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("The other player selected **%s**. Current score is you with %d points and the other player with %d points. Next round! Pedra, papel, or tesoura. Which one you select?", opponentMove, yourScore, otherScore),
		},
	}
}
