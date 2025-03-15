package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    TelegramToken string
    OpenRouterToken string
}

func LoadConfig() (*Config, error) {
    // Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	c := &Config{
        TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		OpenRouterToken: os.Getenv("OPENROUTER_API_TOKEN"),
    }

	if c.TelegramToken == "" {
        return nil,  fmt.Errorf("TELEGRAM_BOT_TOKEN not set")
    }
    if c.OpenRouterToken == "" {
        return nil,  fmt.Errorf("OPENROUTER_API_TOKEN not set")
    }

	return c, nil
}
