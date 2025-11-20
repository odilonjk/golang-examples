package llm

import (
	"context"

	"github.com/odilonjk/pedra-papel-tesoura-llm/internal/models"
)

type LLM interface {
	GetNextMove(ctx context.Context, msgs []models.Message) (string, error)
}
