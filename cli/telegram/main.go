// Package main implements a Telegram bot application that listens for updates
// and handles them concurrently. The bot requires a Telegram API token provided
// through the TELEGRAM_BOT_TOKEN environment variable.
//
// The application features:
// - Graceful shutdown on SIGINT or SIGTERM signals
// - Concurrent handling of incoming updates
// - Configurable update timeout (currently set to 60 seconds)
//
// Environment Variables:
//   - TELEGRAM_BOT_TOKEN: Required. The authentication token for the Telegram Bot API.
//
// The main package initializes the bot through the bot package and starts
// listening for updates. Each update is handled in a separate goroutine.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/olegshulyakov/go-briefly-bot/cli/telegram/bot"
)

// Client timeout.
const timeout = 60

// main is the entry point for the application. It performs the following steps:
// 1. Validates the presence of the TELEGRAM_BOT_TOKEN environment variable
// 2. Initializes a new Telegram bot instance
// 3. Sets up an update channel with a 60-second timeout
// 4. Starts a loop to handle incoming updates concurrently
// 5. Implements graceful shutdown on receiving termination signals
//
// The function exits with status code 1 if:
// - TELEGRAM_BOT_TOKEN is not set
// - Bot initialization fails
//
// The function runs indefinitely until interrupted by SIGINT or SIGTERM signals.
func main() {
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken == "" {
		fmt.Fprintf(os.Stderr, "TELEGRAM_BOT_TOKEN not set")
		os.Exit(1)
	}

	// Start the Telegram bot
	tgBot, err := bot.New(telegramToken)
	if err != nil {
		slog.Error("Failed to create bot", "error", err)
		os.Exit(1)
	}

	// Bot.Debug = true
	slog.Info("Bot started successfully")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout
	updates := tgBot.GetUpdatesChan(u)

	// Handle incoming updates
	for update := range updates {
		go bot.Handle(update)
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down bot...")
}
