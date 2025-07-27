// Package summarization provides functionality to generate text summaries using OpenAI's API.
// It handles localization of prompts, API communication, and error handling.
//
// The package requires the following environment variables to be set:
//   - OPENAI_BASE_URL: The base URL for OpenAI API
//   - OPENAI_API_KEY: The API key for authentication
//   - OPENAI_MODEL: The model to use for summarization
//
// Example usage:
//
//	summary, err := summarization.SummarizeText("Long text to summarize...", "en")
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(summary)
package summarization

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	briefly "github.com/olegshulyakov/go-briefly-bot/lib"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	maxRetries = 3                 // Maximum number of retry attempts for API calls
	maxTimeout = 300 * time.Second // Maximum duration for API requests
)

var (
	openAiBaseURL string // Base URL for OpenAI API
	openAiAPIKey  string // API key for OpenAI authentication
	openAiModel   string // Model identifier for summarization
)

// init initializes the package by loading required environment variables.
// It exits the program if any required variables are missing.
func init() {
	openAiBaseURL = os.Getenv("OPENAI_BASE_URL")
	openAiAPIKey = os.Getenv("OPENAI_API_KEY")
	openAiModel = os.Getenv("OPENAI_MODEL")

	// Validate provider-specific fields
	isError := false
	if openAiBaseURL == "" {
		fmt.Fprintf(os.Stderr, "OPENAI_BASE_URL not set")
		isError = true
	}
	if openAiAPIKey == "" {
		fmt.Fprintf(os.Stderr, "OPENAI_API_KEY not set")
		isError = true
	}
	if openAiModel == "" {
		fmt.Fprintf(os.Stderr, "OPENAI_MODEL not set")
		isError = true
	}

	if isError {
		os.Exit(1)
	}
}

// SummarizeText generates a summary of the input text in the specified language.
// It handles prompt localization, API communication, and response processing.
//
// Parameters:
//   - text: The text content to be summarized
//   - lang: The target language for the summary (e.g., "en", "fr")
//
// Returns:
//   - The generated summary as a string
//   - An error if the summarization fails
//
// The function logs debug information about the summarization process and
// errors if they occur during API communication.
func SummarizeText(text string, lang string) (string, error) {
	slog.Debug("SummarizeText start", "language", lang, "api", openAiBaseURL, "model", openAiModel)

	client := openai.NewClient(
		option.WithBaseURL(openAiBaseURL),
		option.WithAPIKey(openAiAPIKey),
	)

	// Localize system and user prompts
	slog.Debug("Localizing prompts...")
	body := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(briefly.MustLocalizeTemplate(
				lang,
				"llm.prompt",
				map[string]string{
					"text": text,
				},
			)),
		},
		Model: openAiModel,
	}

	slog.Debug("Summarizing text...")
	chatCompletion, err := client.Chat.Completions.New(
		context.Background(),
		body,
		option.WithRequestTimeout(maxTimeout),
		option.WithMaxRetries(maxRetries),
	)

	// Check for errors in the response
	if err != nil {
		slog.Error("Open AI API error", "error", err)
		return "", err
	}

	// Extract summary from response
	slog.Debug("Extracting summary from response...")
	choices := chatCompletion.Choices
	if len(choices) == 0 {
		slog.Warn("Invalid or empty choices in API response, retrying...", "chatCompletion", chatCompletion)
		return "", err
	}

	summary := choices[0].Message.Content

	slog.Debug("SummarizeText completed", "language", lang)
	return summary, nil
}
