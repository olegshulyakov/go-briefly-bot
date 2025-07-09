package main

import (
	"os"

	"github.com/olegshulyakov/go-briefly-bot/briefly"
	"github.com/olegshulyakov/go-briefly-bot/clients/telegram/bot"
)

// main is the entry point for the application.
func main() {
	// Load configuration
	cfg, err := briefly.LoadConfiguration()
	if err != nil {
		briefly.Error("Failed to load config: %v", "error", err)
		os.Exit(1)
	}

	// Start the Telegram bot
	bot.StartTelegramBot(cfg.TelegramToken)
}
