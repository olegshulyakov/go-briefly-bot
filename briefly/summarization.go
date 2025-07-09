package briefly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

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
	Debug("SummarizeText start", "language", lang, "api", Configuration.OpenAiBaseURL, "model", Configuration.OpenAiModel)

	// Localize system and user prompts
	Debug("Localizing prompts...")
	localizer := GetLocalizer(lang)

	systemPrompt := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "llm.system_prompt"})
	userPrompt := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "llm.user_prompt",
		TemplateData: map[string]interface{}{
			"text": text,
		},
	})

	// Prepare payload
	Debug("Preparing API payload...")
	payload := map[string]interface{}{
		"model": Configuration.OpenAiModel,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		Error("Failed to marshal payload", "error", err)
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Determine API endpoint
	var endpoint = "/chat/completions"

	// Retry mechanism
	maxRetries := 3               // Maximum number of retries
	retryDelay := 2 * time.Second // Delay between retries

	for retry := 0; retry < maxRetries; retry++ {
		if retry > 1 {
			Debug("Sleeping after attempt...", "attempt", retry)
			time.Sleep(retryDelay)
		}

		Debug("Attempt to summarize text...", "attempt", retry+1)

		req, err := http.NewRequest("POST", Configuration.OpenAiBaseURL+endpoint, bytes.NewBuffer(payloadBytes))
		if err != nil {
			Error("Failed to create request", "error", err)
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+Configuration.OpenAiAPIKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			Error("Failed to send request", "error", err)
			return "", fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		// Extract summary from response
		Debug("Extracting summary from response...")

		// Decode API response
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			Error("Failed to decode response", "error", err)
			return "", fmt.Errorf("failed to decode response: %w", err)
		}

		// Validate the response structure
		if result == nil {
			Warn("Empty API response, retrying...")
			continue
		}

		// Check for errors in the response
		if errMsg, ok := result["error"].(string); ok {
			Error("API error", "error", errMsg)
			continue
		}

		choices, ok := result["choices"].([]interface{})
		if !ok || len(choices) == 0 {
			Warn("Invalid or empty choices in API response, retrying...", "result", result)
			continue
		}

		firstChoice, ok := choices[0].(map[string]interface{})
		if !ok {
			Warn("Invalid choice structure in API response, retrying...", "choices", choices)
			continue
		}

		message, ok := firstChoice["message"].(map[string]interface{})
		if !ok {
			Warn("Invalid message structure in API response, retrying...", "firstChoice", firstChoice)
			continue
		}

		summary, ok := message["content"].(string)
		if !ok {
			Warn("Invalid content in API response, retrying...")
			continue
		}

		Debug("SummarizeText completed", "language", lang)
		return summary, nil
	}

	// If all retries fail
	return "", fmt.Errorf("failed to summarize text after %d retries", maxRetries)
}
