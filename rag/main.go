package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/odilonjk/golang-examples/rag/openai"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
)

type trail struct {
	Get struct {
		Trail []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"Trail"`
	} `json:"Get"`
}

func main() {
	userQuery := "Trails near Sierra Nevada"

	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	ctx := context.Background()
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error initializing Weaviate. Err: %s", err.Error())
	}

	// Generates embeddings in the DB using OpenAI API.
	// err := seed.Generate(ctx, client)
	// if err != nil {
	// 	log.Fatalf("Error generating seed. Err: %s", err.Error())
	// }

	embedding, err := openai.GetEmbedding(userQuery)
	if err != nil {
		log.Fatalf("Error getting embedding. Err: %s", err.Error())
	}

	nearVector := client.GraphQL().NearVectorArgBuilder().WithVector(embedding)

	res, err := client.GraphQL().Get().
		WithClassName("Trail").
		WithFields(
			graphql.Field{Name: "name"},
			graphql.Field{Name: "description"},
		).
		WithLimit(1). // I could have handled multiple results. But let's keep it simple for now.
		WithNearVector(nearVector).
		Do(ctx)
	if err != nil {
		log.Fatalf("Error fetching Weaviate. Err: %s", err.Error())
	}

	if len(res.Errors) > 0 {
		for _, e := range res.Errors {
			fmt.Errorf("Error fetching data: %s\n", e.Message)
		}
		return
	}

	var trail trail
	b, err := json.Marshal(res.Data)
	if err != nil {
		log.Fatalf("Error parsing into JSON bytes. Err: %s", err.Error())
	}

	err = json.Unmarshal(b, &trail)
	if err != nil {
		log.Fatalf("Error parsing into trail object. Err: %s", err.Error())
	}

	nearTrail := trail.Get.Trail[0]
	trailPrompt := fmt.Sprintf("Nearest trail found: %s — %s", nearTrail.Name, nearTrail.Description)
	msg, err := openai.GetCompletion(ctx, userQuery, trailPrompt)
	if err != nil {
		log.Fatalf("Error on OpenAI chat completion. Err: %s", err.Error())
	}
	log.Printf("[Assistant] %s", msg)
}
