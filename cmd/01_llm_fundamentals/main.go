package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func main() {
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY is not set")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()
	response, err := client.Responses.New(
		ctx,
		responses.ResponseNewParams{
			Model: "gpt-5-nano",
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String("Explain goroutines in one sentence."),
			},
		},
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("LLM request timed out")
			return
		}

		log.Printf("LLM request failed: %v\n", err)
		return
	}

	fmt.Println(response.OutputText())
}
