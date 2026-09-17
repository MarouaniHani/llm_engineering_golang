package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

var systemPrompt = `
You are a snarky assistant that analyzes the contents of a website,
and provides a short, snarky, humorous summary, ignoring text that might be navigation related.
Respond in markdown. Do not wrap the markdown in a code block - respond just with the markdown.
`
var userPromptPrefix = `
Here are the contents of a website.
Provide a short summary of this website.
If it includes news or announcements, then summarize these too.

`

func ScrapeWebsite(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
	}

	// Get title
	title := strings.TrimSpace(doc.Find("title").First().Text())
	if title == "" {
		title = "No title found"
	}

	// Remove irrelevant elements
	doc.Find("script, style, img, input").Remove()

	// Get body text
	text := strings.TrimSpace(doc.Find("body").Text())

	result := title + "\n\n" + text

	// Limit to 2000 characters
	if len(result) > 2000 {
		result = result[:2000]
	}

	return result, nil
}
func messagesFor(website string) []openai.ChatCompletionMessageParamUnion {
	return []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage(userPromptPrefix + website),
	}
}

func summarizeWebsite(client openai.Client, url string) (string, error) {
	website, err := ScrapeWebsite(url)
	if err != nil {
		return "", err
	}
	response, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model:    "gpt-4.1-nano",
			Messages: messagesFor(website),
		},
	)
	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	url := "https://edwarddonner.com"

	summary, err := summarizeWebsite(client, url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(summary)
}
