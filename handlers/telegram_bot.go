package handlers

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"youtube-retell-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

func StartBot(token string) {
	var err error
	Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	Bot.Debug = true

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
	log.Println("Shutting down bot...")
}

func HandleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		switch message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome! Send me a YouTube link, and I'll summarize it for you.")
			Bot.Send(msg)
		}
		return
	}

	// Check if the message contains a YouTube link
	if strings.Contains(message.Text, "youtube.com") || strings.Contains(message.Text, "youtu.be") {
		videoURL := message.Text

		// Fetch transcript
		transcript, err := services.GetTranscript(videoURL)
		if err != nil {
			log.Printf("Failed to get transcript: %v", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Sorry, I couldn't fetch the transcript for this video.")
			Bot.Send(msg)
			return
		}

		// Summarize transcript
		summary, err := services.Summarize(transcript)
		if err != nil {
			log.Printf("Failed to summarize transcript: %v", err)
			msg := tgbotapi.NewMessage(message.Chat.ID, "Sorry, I couldn't summarize the transcript.")
			Bot.Send(msg)
			return
		}

		// Send summary to user
		msg := tgbotapi.NewMessage(message.Chat.ID, summary)
		Bot.Send(msg)
	}
}