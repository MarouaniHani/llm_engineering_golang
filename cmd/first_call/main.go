package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	verified, err := verifyAPIKey(apiKey)
	if err != nil {
		log.Fatal("Error verifying API key: ", err)
	}
	if !verified {
		log.Fatal("API key is not valid")
	}
	message := "Hello, GPT! This is my first ever message to you! Hi!"

	// Create OpenAI client
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	// Call openai api
	response, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model: "gpt-5-nano",
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(message),
			},
		},
	)

	if err != nil {
		fmt.Println("OpenAI error:", err)
		return
	}

	// response.choices[0].message.content
	fmt.Println(response.Choices[0].Message.Content)
}

func verifyAPIKey(apiKey string) (bool, error) {
	if apiKey == "" {
		return false, errors.New("API key is not set")
	}
	if !strings.HasPrefix(apiKey, "sk-proj-") {
		return false, errors.New("API key doesn't start with sk-proj-")
	}
	if strings.TrimSpace(apiKey) != apiKey {
		return false, errors.New("API key contains spaces at the beginning or end")
	}
	return true, nil
}
