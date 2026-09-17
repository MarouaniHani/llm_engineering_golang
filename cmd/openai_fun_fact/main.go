package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	response, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model: "gpt-5-nano",
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("Tell me a fun fact"),
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Choices[0].Message.Content)
}
