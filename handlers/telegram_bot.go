package handlers

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"youtube-retell-bot/config"
	"youtube-retell-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func StartBot(token string) {
	var err error
	Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		config.Logger.Fatalf("Failed to create bot: %v", err)
	}

	Bot.Debug = true
	config.Logger.Info("Bot started successfully")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	// Handle incoming updates
	for update := range updates {
		if update.Message != nil {
			go HandleMessage(update.Message)
		}
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	config.Logger.Info("Shutting down bot...")
}

// sendMessage sends a message using the bot and logs any errors.
func sendMessage(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "Markdown"

	const maxRetries = 3
	var err error

	for i := 0; i < maxRetries; i++ {
		_, err = Bot.Send(msg)
		if err == nil {
			return // Success
		}
		config.Logger.Warnf("Attempt %d: Failed to send message: %v", i+1, err)
	}

	Bot.Send(tgbotapi.NewMessage(chatId, "Sorry, I couldn't summarize the transcript."))

	config.Logger.Errorf("Failed to send message after %d attempts: %v", maxRetries, err)
}

func HandleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		switch message.Command() {
		case "start":
			sendMessage(message.Chat.ID, "Welcome! Send me a YouTube link, and I'll summarize it for you.")
		}
		return
	}

	// Check if the message contains a YouTube link
	if strings.Contains(message.Text, "youtube.com") || strings.Contains(message.Text, "youtu.be") {
		videoURL := message.Text
		config.Logger.Infof("Processing YouTube video: %s", videoURL)

		// Fetch transcript
		transcript, err := services.GetTranscript(videoURL)
		if err != nil {
			config.Logger.Errorf("Failed to get transcript: %v", err)
			sendMessage(message.Chat.ID, "Sorry, I couldn't fetch the transcript for this video.")
			return
		}

		// Summarize transcript
		summary, err := services.Summarize(transcript)
		if err != nil {
			config.Logger.Errorf("Failed to summarize transcript: %v", err)
			sendMessage(message.Chat.ID, "Sorry, I couldn't summarize the transcript.")
			return
		}

		// Send summary to user
		sendMessage(message.Chat.ID, summary)
	}
}