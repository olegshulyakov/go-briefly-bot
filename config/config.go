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
	TelegramToken      string // The token for the Telegram bot.
	YtDlpProxy         string // The proxy URL for yt-dlp.
	LlmProviderType    string // The type of language model provider (e.g., "openai", "ollama").
	SummarizerApiUrl   string // The URL for the OpenAI API.
	SummarizerApiToken string // The token for the OpenAI API.
	SummarizerApiModel string // The model to use for OpenAI.
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
//	cfg, err := loadConfig()
//	if err != nil {
//	    log.Fatalf("Failed to load config: %v", err)
//	}
//	fmt.Println("Telegram Token:", cfg.TelegramToken)
func loadConfig() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		Logger.Warnf("No .env file found, using environment variables.")
	}

	Configuration = &Config{
		TelegramToken:      os.Getenv("TELEGRAM_BOT_TOKEN"),
		YtDlpProxy:         os.Getenv("YT_DLP_PROXY"),
		LlmProviderType:    os.Getenv("LLM_PROVIDER_TYPE"),
		SummarizerApiUrl:   os.Getenv("SUMMARIZER_API_URL"),
		SummarizerApiToken: os.Getenv("SUMMARIZER_API_TOKEN"),
		SummarizerApiModel: os.Getenv("SUMMARIZER_MODEL"),
	}

	// Validate required fields
	if Configuration.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN not set")
	}
	if Configuration.LlmProviderType == "" {
		Logger.Warnf("LLM_PROVIDER_TYPE is not set, setting default OpenAi\n")
		Configuration.LlmProviderType = "openai"
	}

	if Configuration.LlmProviderType != "openai" && Configuration.LlmProviderType != "ollama" {
		return nil, fmt.Errorf("LLM_PROVIDER_TYPE is wrong: %v", Configuration.LlmProviderType)
	}

	// Validate provider-specific fields
	if Configuration.SummarizerApiUrl == "" {
		return nil, fmt.Errorf("SUMMARIZER_API_URL not set")
	}
	if Configuration.SummarizerApiToken == "" {
		return nil, fmt.Errorf("SUMMARIZER_API_TOKEN not set")
	}
	if Configuration.SummarizerApiModel == "" {
		return nil, fmt.Errorf("SUMMARIZER_MODEL not set")
	}

	return Configuration, nil
}

// SetupLogger initializes the global logger with a text formatter and sets the log level to Info.
//
// The logger is configured to output to standard output (stdout).
//
// Example:
//
//	setupLogger()
//	Logger.Info("Logger initialized successfully")
func setupLogger() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.TextFormatter{})
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.DebugLevel)
}
