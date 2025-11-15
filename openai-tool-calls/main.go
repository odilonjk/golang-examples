package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type GetMaterialArgs struct {
	BuildingType string  `json:"building_type"`
	AreaSqM      float64 `json:"area_sqm"`
	MainMaterial string  `json:"main_material"`
}

type Material struct {
	Name         string  `json:"name"`
	Amount       float64 `json:"amount"`
	Unit         string  `json:"unit"`
	InStock      bool    `json:"in_stock"`
	DeliveryDays int     `json:"delivery_days"`
}

type GetMaterialResult struct {
	Materials []Material `json:"materials"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(openaiAPIKey),
	)

	userMessage := "I'd like to build a wooden cabin, around 70 square meters."
	log.Printf("User: %s\n", userMessage)

	params := buildBaseParams()
	params.Messages = append(params.Messages, openai.UserMessage(userMessage))

	ctx := context.Background()
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		log.Fatal(err)
	}
	msg := completion.Choices[0].Message
	if len(msg.ToolCalls) == 0 {
		log.Printf("Assistant without tool: %s", msg.Content)
		return
	}

	params.Messages = append(params.Messages, msg.ToParam())

	tc := msg.ToolCalls[0]
	materialsRequiredJSON := executeGetMaterialToolCall(tc)

	params.Messages = append(params.Messages,
		openai.ToolMessage(
			string(materialsRequiredJSON),
			tc.ID,
		),
	)

	completion, err = client.Chat.Completions.New(ctx, params)
	if err != nil {
		log.Fatalf("Failed to call completion with materials result. Err: %v", err)
	}
	log.Printf("Assistant: %s", completion.Choices[0].Message.Content)
}

func buildBaseParams() openai.ChatCompletionNewParams {
	return openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`
				You are an assistant helping buyers to create orders.

				1. Keep the final answer under 100 words.
				2. Present the materials in a bullet list.
				3. Zero delivery days means "available for immediate delivery".
				5. Do NOT mention the internal inventory list or the tool JSON.
				6. Verify if the user wants to put the order in place or change something.
			`),
		},
		Model: openai.ChatModelGPT5Nano,
		Tools: []openai.ChatCompletionToolUnionParam{
			openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
				Name:        "get_material_required",
				Description: openai.String("Get the list of material required for a building."),
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"building_type": map[string]interface{}{
							"type":        "string",
							"description": "Type of building (e.g. cabin, house, shed)",
						},
						"area_sqm": map[string]interface{}{
							"type":        "number",
							"description": "Approximate area in square meters",
						},
						"main_material": map[string]interface{}{
							"type":        "string",
							"description": "Primary construction material. Available options: wood, raw_logs, bricks.",
						},
					},
					"required": []string{"building_type", "area_sqm", "main_material"},
				},
			}),
		},
	}
}

func executeGetMaterialToolCall(tc openai.ChatCompletionMessageToolCallUnion) []byte {
	var args GetMaterialArgs
	errUnmarshal := json.Unmarshal([]byte(tc.Function.Arguments), &args)
	if errUnmarshal != nil {
		log.Fatal(errUnmarshal)
	}

	materialsRequired := getMaterialRequired(args)

	resultJSON, errMarshal := json.Marshal(materialsRequired)
	if errMarshal != nil {
		log.Fatal(errMarshal)
	}
	return resultJSON
}

func getMaterialRequired(args GetMaterialArgs) GetMaterialResult {
	var materials []Material
	switch args.MainMaterial {
	case "raw_logs":
		materials = []Material{
			{Name: "Raw Logs", Amount: 40, Unit: "pieces", InStock: false, DeliveryDays: 6},
			{Name: "Nails Pack w/ 10 units", Amount: 100, Unit: "package", InStock: true, DeliveryDays: 0},
		}
	case "wood":
		materials = []Material{
			{Name: "Wood Planks", Amount: 200, Unit: "pieces", InStock: false, DeliveryDays: 12},
			{Name: "Nails Pack w/ 10 units", Amount: 60, Unit: "package", InStock: true, DeliveryDays: 0},
		}
	case "bricks":
		materials = []Material{
			{Name: "Bricks", Amount: 30000, Unit: "pieces", InStock: true, DeliveryDays: 0},
			{Name: "Cement", Amount: 20, Unit: "bags", InStock: true, DeliveryDays: 0},
			{Name: "Sand", Amount: 500, Unit: "kg", InStock: false, DeliveryDays: 2},
		}
	}
	return GetMaterialResult{materials}
}
