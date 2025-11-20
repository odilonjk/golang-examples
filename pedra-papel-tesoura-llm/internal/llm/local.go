package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/odilonjk/pedra-papel-tesoura-llm/internal/models"
)

type completionsRequest struct {
	Messages    []models.Message `json:"messages"`
	Model       string           `json:"model"`
	Temperature float32          `json:"temperature"`
}

type Local struct {
	client  *http.Client
	baseURL string
}

const (
	applicationJSON     = "application/json"
	chatCompletionsPath = "/v1/chat/completions"
	model               = "Llama 3 8B Instruct"
)

func NewLocal(baseURL string) *Local {
	c := http.DefaultClient
	return &Local{
		client:  c,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (l *Local) GetNextMove(ctx context.Context, messages []models.Message) (string, error) {
	url := fmt.Sprintf("%s%s", l.baseURL, chatCompletionsPath)
	msgs := completionsRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 0.4,
	}
	msgsJSON, err := json.Marshal(msgs)
	if err != nil {
		return "", err
	}
	res, err := l.client.Post(
		url,
		applicationJSON,
		bytes.NewBuffer(msgsJSON),
	)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
