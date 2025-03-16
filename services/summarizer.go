// Package services provides functionality for interacting with external services,
// such as language models, to perform tasks like text summarization.
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		config.Logger.Errorf("Failed to load config: %v", err)
		return "", fmt.Errorf("failed to load config: %v", err)
	}

	// Determine API URL, token, and model based on LLM provider
	config.Logger.Debugf("Using LLM provider: %v", cfg.LlmProviderType)
	var apiUrl, apiToken, model string

	switch cfg.LlmProviderType {
	case "openai":
		apiUrl = cfg.OpenAiUrl
		apiToken = cfg.OpenAiToken
		model = cfg.OpenAiModel
	case "ollama":
		apiUrl = cfg.OllamaUrl
		apiToken = cfg.OllamaToken
		model = cfg.OllamaModel
	default:
		config.Logger.Errorf("Unsupported LLM provider type: %v", cfg.LlmProviderType)
		return "", fmt.Errorf("unsupported LLM provider type: %v", cfg.LlmProviderType)
	}
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
	if cfg.LlmProviderType == "openai" {
		endpoint = "/chat/completions"
	} else {
		endpoint = "/api/chat"
	}

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

	// Decode API response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		config.Logger.Errorf("Failed to decode response: %v", err)
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	// Extract summary from response
	config.Logger.Debug("Extracting summary from response...")
	summary := result["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	config.Logger.Debugf("SummarizeText completed for language: %v", lang)
	return summary, nil
}
