package main

import (
	"fmt"
	"os"

	"github.com/olegshulyakov/go-briefly-bot/cli/telegram/bot"
)

// main is the entry point for the application.
func main() {
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken == "" {
		fmt.Fprintf(os.Stderr, "TELEGRAM_BOT_TOKEN not set")
		os.Exit(1)
	}

	// Start the Telegram bot
	bot.StartTelegramBot(telegramToken)
}
