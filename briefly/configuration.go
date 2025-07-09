package briefly

import (
	"errors"
	"os"
)

// Config represents the application configuration, including settings for Telegram,
// language model providers (e.g., OpenAI, Ollama), and other environment variables.
type Config struct {
	TelegramToken          string // The token for the Telegram bot.
	YtDlpAdditionalOptions string // The proxy URL for yt-dlp.
	OpenAiBaseURL          string // The URL for the OpenAI API.
	OpenAiAPIKey           string // The token for the OpenAI API.
	OpenAiModel            string // The model to use for OpenAI.
}

var Configuration *Config

// LoadConfiguration loads the application configuration from environment variables.
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
		OpenAiBaseURL:          os.Getenv("OPENAI_BASE_URL"),
		OpenAiAPIKey:           os.Getenv("OPENAI_API_KEY"),
		OpenAiModel:            os.Getenv("OPENAI_MODEL"),
	}

	err := validateConfiguration(Configuration)

	return Configuration, err
}

func validateConfiguration(configuration *Config) error {
	// Validate required fields
	if configuration.TelegramToken == "" {
		return errors.New("TELEGRAM_BOT_TOKEN not set")
	}

	// Validate provider-specific fields
	if configuration.OpenAiBaseURL == "" {
		return errors.New("OPENAI_BASE_URL not set")
	}
	if configuration.OpenAiAPIKey == "" {
		return errors.New("OPENAI_API_KEY not set")
	}
	if configuration.OpenAiModel == "" {
		return errors.New("OPENAI_MODEL not set")
	}

	return nil
}
