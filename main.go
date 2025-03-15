package main

import (
	"youtube-retell-bot/config"
	"youtube-retell-bot/handlers"
)

func main() {
    // Set up logger
	config.SetupLogger()

	// Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        config.Logger.Fatalf("Failed to load config: %v", err)
    }

    // Start the Telegram bot
	handlers.StartBot(cfg.TelegramToken)
}