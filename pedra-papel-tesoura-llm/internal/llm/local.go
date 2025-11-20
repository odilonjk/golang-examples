package llm

import (
	"context"
	"net/http"
)

type Local struct {
	client  *http.Client
	baseURL string
}

func NewLocal(baseURL string) *Local {
	c := http.DefaultClient
	return &Local{
		client:  c,
		baseURL: baseURL,
	}
}

func (l *Local) GetNextMove(ctx context.Context) (string, error) {
	return "pedra", nil
}
