package handlers

import (
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	"youtube-retell-bot/config"
	"youtube-retell-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	Bot             *tgbotapi.BotAPI
	userLastRequest = make(map[int64]time.Time) // Tracks the last request time for each user
	userMutex       = sync.Mutex{}              // Mutex to protect userLastRequest map
)

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

// sendWithRetry sends a message with retries.
func sendWithRetry(msg tgbotapi.Chattable) (tgbotapi.Message, error) {
	const maxRetries = 3
	var err error
	var sentMsg tgbotapi.Message

	for i := 0; i < maxRetries; i++ {
		sentMsg, err = Bot.Send(msg)
		if err == nil {
			return sentMsg, nil // Success
		}
		config.Logger.Warnf("Attempt %d: Failed to send message: %v", i+1, err)
	}

	config.Logger.Errorf("Failed to send message after %d attempts: %v", maxRetries, err)
	return tgbotapi.Message{}, err
}

// sendMessage sends a message using the bot and logs any errors.
func sendMessage(userMessage *tgbotapi.Message, text string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(userMessage.Chat.ID, text)
	msg.ReplyToMessageID = userMessage.MessageID
	return sendWithRetry(msg)
}

// editMessage updates an existing message with new text.
func editMessage(userMessage *tgbotapi.Message, messageToEdit tgbotapi.Message, text string) (tgbotapi.Message, error) {
	editMsg := tgbotapi.NewEditMessageText(userMessage.Chat.ID, messageToEdit.MessageID, text)
	return sendWithRetry(editMsg)
}

// sendErrorMessage sends an error message to the user.
func sendErrorMessage(userMessage *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(userMessage.Chat.ID, text)
	msg.ReplyToMessageID = userMessage.MessageID
	Bot.Send(msg)
}

// isUserRateLimited checks if the user has made a request within the last 10 seconds.
func isUserRateLimited(userId int64) bool {
	userMutex.Lock()
	defer userMutex.Unlock()

	lastRequest, exists := userLastRequest[userId]
	if exists && time.Since(lastRequest) < 10*time.Second {
		return true // User is rate-limited
	}

	userLastRequest[userId] = time.Now() // Update the last request time
	return false
}

// handleTelegramMessage processes incoming messages from users.
func handleTelegramMessage(message *tgbotapi.Message) {
	// Determine the user's language (default to English)
	userLanguage := message.From.LanguageCode
	if userLanguage == "" {
		userLanguage = "en"
	}

	// Create a localizer for the user's language
	localizer := config.GetLocalizer(userLanguage)

	// Check if the user is rate-limited
	if isUserRateLimited(message.From.ID) {
		sendErrorMessage(message, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.errors.rate_limited"}))
		return
	}

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			sendMessage(message, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.welcome.message"}))
		}
		return
	}

	// Check if the message contains a YouTube link
	if strings.Contains(message.Text, "youtube.com") || strings.Contains(message.Text, "youtu.be") {
		videoURL := message.Text
		config.Logger.Infof("Processing YouTube video: %s", videoURL)

		processingMsg, err := sendMessage(message, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.progress.processing"}))
		if err != nil {
			config.Logger.Errorf("Failed to send processing message: %v", err)
			return
		}

		// Fetch video info
		processingMsg, err = editMessage(message, processingMsg, processingMsg.Text+"\n"+localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.progress.fetching_info"}))
		if err != nil {
			config.Logger.Errorf("Failed to update progress message: %v", err)
			return
		}

		videoInfo, err := services.GetYoutubeVideoInfo(videoURL)
		if err != nil {
			config.Logger.Errorf("Failed to get video info: %v", err)
			editMessage(message, processingMsg, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.errors.info_failed"}))
			return
		}

		// Fetch transcript
		processingMsg, err = editMessage(message, processingMsg, processingMsg.Text+"\n"+localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.progress.fetching_transcript"}))
		if err != nil {
			config.Logger.Errorf("Failed to update progress message: %v", err)
			return
		}

		transcript, err := services.GetYoutubeTranscript(videoURL)
		if err != nil {
			config.Logger.Errorf("Failed to get transcript: %v", err)
			editMessage(message, processingMsg, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.errors.transcript_failed"}))
			return
		}

		// Summarize transcript
		processingMsg, err = editMessage(message, processingMsg, processingMsg.Text+"\n"+localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.progress.summarizing"}))
		if err != nil {
			config.Logger.Errorf("Failed to update progress message: %v", err)
			return
		}

		summary, err := services.SummarizeText(transcript, userLanguage)
		if err != nil {
			config.Logger.Errorf("Failed to summarize transcript: %v", err)
			editMessage(message, processingMsg, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.errors.summary_failed"}))
			return
		}

		// Send summary to user
		_, err = sendMessage(message, localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "telegram.result.summary",
			TemplateData: map[string]interface{}{
				"title": videoInfo.Title,
				"text":  summary,
			},
		}))
		if err != nil {
			sendErrorMessage(message, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.errors.general"}))
		}

		// Delete the "Processing" message
		deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, processingMsg.MessageID)
		Bot.Send(deleteMsg)
	}
}
