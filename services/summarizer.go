package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"youtube-retell-bot/config"
)

func Summarize(text string) (string, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %v", err)
	}

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
		return "", fmt.Errorf("unsupported LLM provider type: %v", cfg.LlmProviderType)
	}

	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant that retells text."},
			{"role": "user", "content": "Summarize the retell text: " + text},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	var endpoint string
	if cfg.LlmProviderType == "openai" {
		endpoint = "/chat/completions"
	} else {
		endpoint = "/api/chat"
	}

	req, err := http.NewRequest("POST", apiUrl+endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	summary := result["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return summary, nil
}
