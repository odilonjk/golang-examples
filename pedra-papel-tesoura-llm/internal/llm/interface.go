package llm

import "context"

type LLM interface {
	GetNextMove(ctx context.Context) (string, error)
}
