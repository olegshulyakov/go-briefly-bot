package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

type Config struct {
	TelegramToken string

	LlmProviderType string

	OpenAiUrl   string
	OpenAiToken string
	OpenAiModel string

	OllamaUrl   string
	OllamaToken string
	OllamaModel string
}

func LoadConfig() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	c := &Config{
		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),

		LlmProviderType: os.Getenv("LLM_PROVIDER_TYPE"),

		OpenAiUrl:   os.Getenv("OPEN_AI_API_URL"),
		OpenAiToken: os.Getenv("OPEN_AI_API_TOKEN"),
		OpenAiModel: os.Getenv("OPEN_AI_MODEL"),

		OllamaUrl:   os.Getenv("OLLAMA_API_URL"),
		OllamaToken: os.Getenv("OLLAMA_API_TOKEN"),
		OllamaModel: os.Getenv("OLLAMA_MODEL"),
	}

	if c.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN not set")
	}
	if c.LlmProviderType == "" {
		return nil, fmt.Errorf("LLM_PROVIDER_TYPE is not set")
	}

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

func SetupLogger() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.TextFormatter{})
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.InfoLevel)
}
