package main

import (
	"log"

	"youtube-retell-bot/config"
	"youtube-retell-bot/handlers"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Start the Telegram bot
	handlers.StartBot(cfg.TelegramToken)
}