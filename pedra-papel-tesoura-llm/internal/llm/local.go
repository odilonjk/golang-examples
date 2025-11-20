package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type completionsRequest struct {
	Messages    []message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature"`
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

func (l *Local) GetNextMove(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s%s", l.baseURL, chatCompletionsPath)
	msgs := completionsRequest{
		Model: model,
		Messages: []message{
			{Role: "system", Content: "You're playing the game 'Pedra, papel ou tesoura.'. You just need to select one of these values: pedra, or papel, or tesoura. Pedra is stronger than tesoura. Teroura is stronger than papel. Papel is stronger than pedra. You don't know which one your opponent will select, so it's just a matter of being lucky."},
			{Role: "user", Content: "Your opponent is selecting one of the three options. Which one you select?"},
		},
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
