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

// main is the entry point for the application.
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
	u.Timeout = 60
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
