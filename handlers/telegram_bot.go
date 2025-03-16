package handlers

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"youtube-retell-bot/config"
	"youtube-retell-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var Bot *tgbotapi.BotAPI

func StartTelegramBot(token string) {
	var err error
	Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		config.Logger.Fatalf("Failed to create bot: %v", err)
	}

	// Bot.Debug = true
	config.Logger.Info("Bot started successfully")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	// Handle incoming updates
	for update := range updates {
		if update.Message != nil {
			go handleTelegramMessage(update.Message)
		}
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	config.Logger.Info("Shutting down bot...")
}

// sendMessage sends a message using the bot and logs any errors.
func sendMessage(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)

	const maxRetries = 3
	var err error

	for i := 0; i < maxRetries; i++ {
		_, err = Bot.Send(msg)
		if err == nil {
			return nil // Success
		}
		config.Logger.Warnf("Attempt %d: Failed to send message: %v", i+1, err)
	}

	config.Logger.Errorf("Failed to send message after %d attempts: %v", maxRetries, err)

	return err
}

func sendErrorMessage(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	Bot.Send(msg)
}

func handleTelegramMessage(message *tgbotapi.Message) {
	// Determine the user's language (default to English)
	userLanguage := message.From.LanguageCode
	if userLanguage == "" {
		userLanguage = "en"
	}

	// Create a localizer for the user's language
	localizer := config.GetLocalizer(userLanguage)

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			sendMessage(message.Chat.ID, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.welcome.message"}))
		}
		return
	}

	// Check if the message contains a YouTube link
	if strings.Contains(message.Text, "youtube.com") || strings.Contains(message.Text, "youtu.be") {
		videoURL := message.Text
		config.Logger.Infof("Processing YouTube video: %s", videoURL)

		// Fetch video info
		videoInfo, err := services.GetYoutubeVideoInfo(videoURL)
		if err != nil {
			config.Logger.Errorf("Failed to get video info: %v", err)
			sendErrorMessage(message.Chat.ID, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.errors.info_failed"}))
			return
		}

		// Fetch transcript
		transcript, err := services.GetYoutubeTranscript(videoURL)
		if err != nil {
			config.Logger.Errorf("Failed to get transcript: %v", err)
			sendErrorMessage(message.Chat.ID, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.errors.transcript_failed"}))
			return
		}

		// Summarize transcript
		summary, err := services.SummarizeText(transcript, userLanguage)
		if err != nil {
			config.Logger.Errorf("Failed to summarize transcript: %v", err)
			sendErrorMessage(message.Chat.ID, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.errors.summary_failed"}))
			return
		}

		// Send summary to user
		err = sendMessage(message.Chat.ID, localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "telegram.result.summary",
			TemplateData: map[string]interface{}{
				"title": videoInfo.Title,
				"text":  summary,
			},
		}))
		if err != nil {
			sendErrorMessage(message.Chat.ID, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.errors.general"}))
		}
	}
}
