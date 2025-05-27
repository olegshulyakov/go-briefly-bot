// Package main is the entry point for the YouTube Retell Bot application.
//
// The application initializes the logger, localizer, and configuration, then starts
// the Telegram bot using the provided token.
package main

import (
	"youtube-briefly-bot/config"
	"youtube-briefly-bot/handlers"
)

// main is the entry point for the application.
//
// It performs the following steps:
//  1. Initializes the logger.
//  2. Initializes the localizer.
//  3. Loads the application configuration.
//  4. Starts the Telegram bot using the configured token.
//
// Example:
//
//	func main() {
//	    // Set up logger
//	    config.SetupLogger()
//
//	    // Set up localizer
//	    config.SetupLocalizer()
//
//	    // Load configuration
//	    cfg, err := config.LoadConfig()
//	    if err != nil {
//	        config.Logger.Fatalf("Failed to load config: %v", err)
//	    }
//
//	    // Start the Telegram bot
//	    handlers.StartTelegramBot(cfg.TelegramToken)
//	}
func main() {
	// Load configuration
	cfg, err := config.SetupConfig()
	if err != nil {
		config.Logger.Fatalf("Failed to load config: %v", err)
	}

	// Start the Telegram bot
	handlers.StartTelegramBot(cfg.TelegramToken)
}
