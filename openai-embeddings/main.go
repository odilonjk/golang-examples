package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(openaiAPIKey),
	)
	res, err := client.Embeddings.New(context.Background(), openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String("The Rubicon Trail is a 22-mile-long route, part road and part 4x4 trail, located in the Sierra Nevada of the western United States, due west of Lake Tahoe and about 80 miles (130 km) east of Sacramento."),
		},
		Model: "text-embedding-3-small",
	})
	if err != nil {
		log.Fatalf("Error when generating embeddings. Err: %s", err.Error())
	}
	log.Printf("Generated embeddings: %v", res.Data[0].Embedding)
}
