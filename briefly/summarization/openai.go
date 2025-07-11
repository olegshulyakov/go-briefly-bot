package summarization

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/olegshulyakov/go-briefly-bot/briefly"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// maxRetries sets maximum number of retries
const maxRetries = 3

// maxTimeout sets maximun request timeout
const maxTimeout = 20 * time.Second

var (
	openAiBaseURL string
	openAiAPIKey  string
	openAiModel   string
)

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

// SummarizeText sends a request to a configured Language Model (LLM) provider
// (e.g., OpenAI or Ollama) to summarize the given text in the specified language.
//
// The function performs the following steps:
//  1. Loads the application configuration to determine the LLM provider and its settings.
//  2. Localizes the system and user prompts based on the specified language.
//  3. Prepares and sends an HTTP request to the LLM provider's API.
//  4. Decodes the API response and extracts the summarized text.
//
// Parameters:
//   - text: The text to be summarized.
//   - lang: The language code (e.g., "en", "ru") for localization and summarization.
//
// Returns:
//   - A string containing the summarized text.
//   - An error if any step fails, such as configuration loading, API request, or response decoding.
//
// Example:
//
//	summary, err := SummarizeText("This is a long text to summarize.", "en")
//	if err != nil {
//	    log.Errorf("Failed to summarize text: %v", err)
//	}
//	fmt.Println("Summary:", summary)
//
// Notes:
//   - The function relies on the application configuration (`LoadConfig`) to determine
//     the LLM provider (e.g., OpenAI or Ollama) and its settings (e.g., API URL, token, model).
//   - The function uses the `go-i18n` package for localization of system and user prompts.
//   - The API response is expected to contain a "choices" field with the summarized text.
func SummarizeText(text string, lang string) (string, error) {
	briefly.Debug("SummarizeText start", "language", lang, "api", openAiBaseURL, "model", openAiModel)

	// Localize system and user prompts
	briefly.Debug("Localizing prompts...")
	localizer := briefly.GetLocalizer(lang)

	client := openai.NewClient(
		option.WithBaseURL(openAiBaseURL),
		option.WithAPIKey(openAiAPIKey),
	)

	body := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "llm.system_prompt"})),
			openai.UserMessage(text),
		},
		Model: openAiModel,
	}

	briefly.Debug("Summarizing text...")

	chatCompletion, err := client.Chat.Completions.New(
		context.Background(),
		body,
		option.WithRequestTimeout(maxTimeout),
		option.WithMaxRetries(maxRetries),
	)

	// Check for errors in the response
	if err != nil {
		briefly.Error("Open AI API error", "error", err)
		return "", err
	}

	// Extract summary from response
	briefly.Debug("Extracting summary from response...")

	choices := chatCompletion.Choices
	if len(choices) == 0 {
		briefly.Warn("Invalid or empty choices in API response, retrying...", "chatCompletion", chatCompletion)
		return "", err
	}

	summary := choices[0].Message.Content

	briefly.Debug("SummarizeText completed", "language", lang)
	return summary, nil
}
