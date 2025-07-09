package briefly

import (
	"fmt"
	"os"
)

// Config represents the application configuration, including settings for Telegram,
// language model providers (e.g., OpenAI, Ollama), and other environment variables.
type Config struct {
	TelegramToken          string // The token for the Telegram bot.
	YtDlpAdditionalOptions string // The proxy URL for yt-dlp.
	OpenAiBaseUrl          string // The URL for the OpenAI API.
	OpenAiApiKey           string // The token for the OpenAI API.
	OpenAiModel            string // The model to use for OpenAI.
}

var Configuration *Config

// LoadConfig loads the application configuration from environment variables.
//
// The function reads the `.env` file and populates the Config struct with the values.
// It also performs validation to ensure required fields are set.
//
// Returns:
//   - A pointer to the Config struct containing the loaded configuration.
//   - An error if the `.env` file cannot be loaded or if required fields are missing.
//
// Example:
//
//	cfg, err := LoadConfiguration()
//	if err != nil {
//	    log.Fatalf("Failed to load config: %v", err)
//	}
func LoadConfiguration() (*Config, error) {
	Configuration = &Config{
		TelegramToken:          os.Getenv("TELEGRAM_BOT_TOKEN"),
		YtDlpAdditionalOptions: os.Getenv("YT_DLP_ADDITIONAL_OPTIONS"),
		OpenAiBaseUrl:          os.Getenv("OPENAI_BASE_URL"),
		OpenAiApiKey:           os.Getenv("OPENAI_API_KEY"),
		OpenAiModel:            os.Getenv("OPENAI_MODEL"),
	}

	// Validate required fields
	if Configuration.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN not set")
	}

	// Validate provider-specific fields
	if Configuration.OpenAiBaseUrl == "" {
		return nil, fmt.Errorf("OPENAI_BASE_URL not set")
	}
	if Configuration.OpenAiApiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}
	if Configuration.OpenAiModel == "" {
		return nil, fmt.Errorf("OPENAI_MODEL not set")
	}

	return Configuration, nil
}
