// Package config provides functionality for loading and managing application configuration,
// including environment variables, logging, and localization.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Logger is a global logger instance used for logging throughout the application.
var Logger *logrus.Logger

// Config represents the application configuration, including settings for Telegram,
// language model providers (e.g., OpenAI, Ollama), and other environment variables.
type Config struct {
	TelegramToken   string // The token for the Telegram bot.
	LlmProviderType string // The type of language model provider (e.g., "openai", "ollama").
	OpenAiUrl       string // The URL for the OpenAI API.
	OpenAiToken     string // The token for the OpenAI API.
	OpenAiModel     string // The model to use for OpenAI.
	OllamaUrl       string // The URL for the Ollama API.
	OllamaToken     string // The token for the Ollama API.
	OllamaModel     string // The model to use for Ollama.
}

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
//	cfg, err := LoadConfig()
//	if err != nil {
//	    log.Fatalf("Failed to load config: %v", err)
//	}
//	fmt.Println("Telegram Token:", cfg.TelegramToken)
func LoadConfig() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	c := &Config{
		TelegramToken:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		LlmProviderType: os.Getenv("LLM_PROVIDER_TYPE"),
		OpenAiUrl:       os.Getenv("OPEN_AI_API_URL"),
		OpenAiToken:     os.Getenv("OPEN_AI_API_TOKEN"),
		OpenAiModel:     os.Getenv("OPEN_AI_MODEL"),
		OllamaUrl:       os.Getenv("OLLAMA_API_URL"),
		OllamaToken:     os.Getenv("OLLAMA_API_TOKEN"),
		OllamaModel:     os.Getenv("OLLAMA_MODEL"),
	}

	// Validate required fields
	if c.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN not set")
	}
	if c.LlmProviderType == "" {
		return nil, fmt.Errorf("LLM_PROVIDER_TYPE is not set")
	}

	// Validate provider-specific fields
	switch c.LlmProviderType {
	case "openai":
		if c.OpenAiUrl == "" {
			return nil, fmt.Errorf("OPEN_AI_API_URL not set")
		}
		if c.OpenAiToken == "" {
			return nil, fmt.Errorf("OPEN_AI_API_TOKEN not set")
		}
		if c.OpenAiModel == "" {
			return nil, fmt.Errorf("OPEN_AI_MODEL not set")
		}
	case "ollama":
		if c.OllamaUrl == "" {
			return nil, fmt.Errorf("OLLAMA_API_URL not set")
		}
		if c.OllamaToken == "" {
			return nil, fmt.Errorf("OLLAMA_API_TOKEN not set")
		}
		if c.OllamaModel == "" {
			return nil, fmt.Errorf("OLLAMA_MODEL not set")
		}
	default:
		return nil, fmt.Errorf("LLM_PROVIDER_TYPE is wrong: %v", c.LlmProviderType)
	}

	return c, nil
}

// SetupLogger initializes the global logger with a text formatter and sets the log level to Info.
//
// The logger is configured to output to standard output (stdout).
//
// Example:
//
//	SetupLogger()
//	Logger.Info("Logger initialized successfully")
func SetupLogger() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.TextFormatter{})
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.InfoLevel)
}
