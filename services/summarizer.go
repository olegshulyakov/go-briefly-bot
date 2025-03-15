package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"youtube-retell-bot/config"
)

func Summarize(text string) (string, error) {
	// Use OpenRouter API for summarization
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %v", err)
	}

	payload := map[string]interface{}{
		"model": "mistralai/mistral-nemo:free",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant that retells text."},
			{"role": "user", "content": "Summarize the retell text: " + text},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.OpenRouterToken)
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