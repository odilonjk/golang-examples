package main

import (
	"github.com/odilonjk/pedra-papel-tesoura-llm/internal/game"
	"github.com/odilonjk/pedra-papel-tesoura-llm/internal/llm"
)

func main() {
	llmClient := llm.NewLocal("localhost:4891")
	game := game.New(llmClient)
	game.Start()
}
