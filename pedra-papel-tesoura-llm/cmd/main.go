package main

import (
	"context"

	"github.com/odilonjk/pedra-papel-tesoura-llm/internal/game"
	"github.com/odilonjk/pedra-papel-tesoura-llm/internal/llm"
)

func main() {
	ctx := context.Background()
	llmClient := llm.NewLocal("localhost:4891")
	game := game.New(llmClient)
	game.Start(ctx)
}
