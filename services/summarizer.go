// Package services provides functionality for interacting with external services,
// such as language models, to perform tasks like text summarization.
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"youtube-retell-bot/config"

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
//   - The function relies on the application configuration (`config.LoadConfig`) to determine
//     the LLM provider (e.g., OpenAI or Ollama) and its settings (e.g., API URL, token, model).
//   - The function uses the `go-i18n` package for localization of system and user prompts.
//   - The API response is expected to contain a "choices" field with the summarized text.
//   - Logging is performed using the `config.Logger` for debugging and error tracking.
func SummarizeText(text string, lang string) (string, error) {
	config.Logger.Debugf("SummarizeText start for language: %v", lang)

	// Determine API URL, token, and model based on LLM provider
	if config.Configuration.LlmProviderType != "openai" && config.Configuration.LlmProviderType != "ollama" {
		config.Logger.Errorf("Unsupported LLM provider type: %v", config.Configuration.LlmProviderType)
		return "", fmt.Errorf("unsupported LLM provider type: %v", config.Configuration.LlmProviderType)
	}

	config.Logger.Debugf("Using LLM provider: %v", config.Configuration.LlmProviderType)

	apiUrl := config.Configuration.SummarizerApiUrl
	apiToken := config.Configuration.SummarizerApiToken
	model := config.Configuration.SummarizerApiModel

	config.Logger.Debugf("API URL: %v, Model: %v", apiUrl, model)

	// Localize system and user prompts
	config.Logger.Debug("Localizing prompts...")
	localizer := config.GetLocalizer(lang)

	systemPrompt := localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "llm.system_prompt"})
	userPrompt := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "llm.user_prompt",
		TemplateData: map[string]interface{}{
			"text": text,
		},
	})

	// Prepare payload
	config.Logger.Debug("Preparing API payload...")
	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		config.Logger.Errorf("Failed to marshal payload: %v", err)
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Determine API endpoint
	var endpoint string
	if config.Configuration.LlmProviderType == "openai" {
		endpoint = "/chat/completions"
	} else {
		endpoint = "/api/chat"
	}

	// Retry mechanism
	maxRetries := 10              // Maximum number of retries
	retryDelay := 2 * time.Second // Delay between retries

	for retry := 0; retry < maxRetries; retry++ {
		if retry > 1 {
			config.Logger.Debugf("Sleeping after %v attempt...", retry)
			time.Sleep(retryDelay)
		}

		config.Logger.Debugf("Attempt %d to summarize text...", retry+1)

		req, err := http.NewRequest("POST", apiUrl+endpoint, bytes.NewBuffer(payloadBytes))
		if err != nil {
			config.Logger.Errorf("Failed to create request: %v", err)
			return "", fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+apiToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			config.Logger.Errorf("Failed to send request: %v", err)
			return "", fmt.Errorf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// Extract summary from response
		config.Logger.Debug("Extracting summary from response...")

		// Decode API response
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			config.Logger.Errorf("Failed to decode response: %v", err)
			return "", fmt.Errorf("failed to decode response: %v", err)
		}

		// Validate the response structure
		if result == nil {
			config.Logger.Warn("Empty API response, retrying...")
			continue
		}

		// Check for errors in the response
		if errMsg, ok := result["error"].(string); ok {
			config.Logger.Errorf("API error: %v", errMsg)
			continue
		}

		choices, ok := result["choices"].([]interface{})
		if !ok || len(choices) == 0 {
			config.Logger.Warnf("Invalid or empty choices in API response, retrying...\n%v", result)
			continue
		}

		firstChoice, ok := choices[0].(map[string]interface{})
		if !ok {
			config.Logger.Warnf("Invalid choice structure in API response, retrying...\n%v", choices)
			continue
		}

		message, ok := firstChoice["message"].(map[string]interface{})
		if !ok {
			config.Logger.Warnf("Invalid message structure in API response, retrying...\n%v", firstChoice)
			continue
		}

		summary, ok := message["content"].(string)
		if !ok {
			config.Logger.Warn("Invalid content in API response, retrying...")
			continue
		}

		config.Logger.Debugf("SummarizeText completed for language: %v", lang)
		return summary, nil
	}

	// If all retries fail
	return "", fmt.Errorf("failed to summarize text after %d retries", maxRetries)
}
