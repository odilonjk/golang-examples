package openai

// I could have handled better the client usage to avoid generating it multiple times
// Same for the API key.

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func GetEmbedding(input string) ([]float32, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(openaiAPIKey),
	)
	res, err := client.Embeddings.New(context.Background(), openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(input),
		},
		Model: "text-embedding-3-small",
	})
	if err != nil {
		return nil, err
	}
	embedding := res.Data[0].Embedding
	embedding32 := make([]float32, len(embedding))
	for i, v := range embedding {
		embedding32[i] = float32(v)
	}
	return embedding32, nil
}

func GetCompletion(ctx context.Context, userQuery, trailPrompt string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(openaiAPIKey),
	)
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`
				You are an assistant helping drivers to find trails.

				1. Keep the final answer under 100 words.
				2. Present the nearest trail found in the DB.
				3. Verify if the user wants to send it to the GPS.
			`),
		},
		Model: openai.ChatModelGPT5Nano,
	}
	params.Messages = append(params.Messages, openai.UserMessage(userQuery))
	params.Messages = append(params.Messages, openai.SystemMessage(trailPrompt))
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", err
	}
	return completion.Choices[0].Message.Content, nil
}
