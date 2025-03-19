// Package handlers provides functionality to handle Telegram bot interactions,
// including message processing, rate limiting, and YouTube video summarization.
package handlers

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"youtube-retell-bot/config"
	"youtube-retell-bot/services"
	"youtube-retell-bot/utils"

	tg_md2html "github.com/PaulSonOfLars/gotg_md2html"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	// Bot is the global Telegram bot instance.
	Bot *tgbotapi.BotAPI

	// userLastRequest tracks the last request time for each user to enforce rate limiting.
	userLastRequest = make(map[int64]time.Time)

	// userMutex is a mutex to protect concurrent access to the userLastRequest map.
	userMutex = sync.Mutex{}
)

// StartTelegramBot initializes and starts the Telegram bot with the provided token.
// It handles incoming updates and gracefully shuts down on receiving termination signals.
//
// Parameters:
//   - token: The Telegram bot token used for authentication.
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

// sendWithRetry sends a message with retries in case of failure.
//
// Parameters:
//   - msg: The message to be sent.
//
// Returns:
//   - tgbotapi.Message: The message that was sent.
//   - error: An error if the message fails to send after retries.
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

// sendMessage sends a message to the specified chat and logs any errors.
//
// Parameters:
//   - userMessage: The original message to reply to.
//   - text: The text of the message to send.
//
// Returns:
//   - tgbotapi.Message: The message that was sent.
//   - error: An error if the message fails to send.
func sendMessage(userMessage *tgbotapi.Message, text string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(userMessage.Chat.ID, text)
	msg.ReplyToMessageID = userMessage.MessageID
	return sendWithRetry(msg)
}

// sendMarkdownMessage sends a message with Markdown formatting to the specified chat.
//
// Parameters:
//   - userMessage: The original message to reply to.
//   - markdownText: The message text formatted in Markdown.
//
// Returns:
//   - tgbotapi.Message: The message that was sent.
//   - error: An error if the message fails to send.
func sendMarkdownMessage(userMessage *tgbotapi.Message, text string) (tgbotapi.Message, error) {
	escapedMessageText := tg_md2html.MD2HTMLV2(text)

	msg := tgbotapi.NewMessage(userMessage.Chat.ID, escapedMessageText)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true
	return sendWithRetry(msg)
}

// editMessage updates an existing message with new text.
//
// Parameters:
//   - userMessage: The original message to reply to.
//   - messageToEdit: The message to be edited.
//   - text: The new text for the message.
//
// Returns:
//   - tgbotapi.Message: The edited message.
//   - error: An error if the message fails to update.
func editMessage(userMessage *tgbotapi.Message, messageToEdit tgbotapi.Message, text string) (tgbotapi.Message, error) {
	editMsg := tgbotapi.NewEditMessageText(userMessage.Chat.ID, messageToEdit.MessageID, text)
	return sendWithRetry(editMsg)
}

// sendErrorMessage sends an error message to the user.
//
// Parameters:
//   - userMessage: The original message to reply to.
//   - text: The error message text.
func sendErrorMessage(userMessage *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(userMessage.Chat.ID, text)
	msg.ReplyToMessageID = userMessage.MessageID
	Bot.Send(msg)
}

// isUserRateLimited checks if the user has made a request within the last 30 seconds.
//
// Parameters:
//   - userId: The ID of the user to check.
//
// Returns:
//   - bool: True if the user is rate-limited, false otherwise.
func isUserRateLimited(userId int64) bool {
	userMutex.Lock()
	defer userMutex.Unlock()

	lastRequest, exists := userLastRequest[userId]
	if exists && time.Since(lastRequest) < 30*time.Second {
		return true // User is rate-limited
	}

	userLastRequest[userId] = time.Now() // Update the last request time
	return false
}

// handleTelegramMessage processes incoming messages from users, including commands and YouTube links.
//
// Parameters:
//   - message: The incoming message to process.
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
		config.Logger.Warnf("Rate Limit exceeded: userId=%v", message.From.ID)
		sendErrorMessage(message, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.error.rate_limited"}))
		return
	}

	text := message.Text
	config.Logger.Debugf("Request: userId=%v, user='%v', text=%s", message.From.ID, message.From, text)

	if message.IsCommand() {
		switch message.Command() {
		case "start":
			sendMessage(message, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.welcome.message"}))
		}
		return
	}

	// Check if the message contains a YouTube link and extract URL
	videoURLs, err := services.ExtractAllYouTubeURLs(text)
	if err != nil {
		config.Logger.Errorf("Got invalid processing message: userId=%v, text=%v", message.From.ID, text)
		sendErrorMessage(message, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.error.no_url_found"}))
		return
	}

	// Notify user about process start
	config.Logger.Infof("Processing YouTube video: %s", videoURLs)
	processingMsg, err := sendMessage(message, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.progress.processing"}))
	if err != nil {
		config.Logger.Errorf("Failed to send processing message: %v", err)
		return
	}

	// Check if there are multiple URLs
	if len(videoURLs) > 1 {
		processingMsg, err = editMessage(message, processingMsg, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.error.multiple_urls"}))
	}
	videoURL := videoURLs[0]

	// Fetch video info
	processingMsg, err = editMessage(message, processingMsg, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.progress.fetching_info"}))
	if err != nil {
		config.Logger.Errorf("Failed to update progress message: %v", err)
		return
	}

	videoInfo, err := services.GetYoutubeVideoInfo(videoURL)
	if err != nil {
		config.Logger.Errorf("Failed to get video info: userId=%v, videoURL=%v, err=%v", message.From.ID, videoURL, err)
		editMessage(message, processingMsg, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.error.info_failed"}))
		return
	}

	// Fetch transcript
	processingMsg, err = editMessage(message, processingMsg, processingMsg.Text+"\n"+localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.progress.fetching_transcript"}))
	if err != nil {
		config.Logger.Errorf("Failed to update progress message: userId=%v, err=%v", message.From.ID, err)
		return
	}

	transcript, err := services.GetYoutubeTranscript(videoURL, userLanguage)
	if err != nil {
		config.Logger.Errorf("Failed to get transcript: userId=%v, videoURL=%v, err=%v", message.From.ID, videoURL, err)
		editMessage(message, processingMsg, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.error.transcript_failed"}))
		return
	}

	// Summarize transcript
	processingMsg, err = editMessage(message, processingMsg, processingMsg.Text+"\n"+localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.progress.summarizing"}))
	if err != nil {
		config.Logger.Errorf("Failed to update progress message: userId=%v, err=%v", message.From.ID, err)
		return
	}

	summary, err := services.SummarizeText(transcript, userLanguage)
	if err != nil {
		config.Logger.Errorf("Failed to summarize transcript: userId=%v, videoURL=%v, err=%v", message.From.ID, videoURL, err)
		editMessage(message, processingMsg, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.error.summary_failed"}))
		return
	}

	// Send summary to user
	chunkSize := 4000 - len(videoInfo.Title)
	chunks := utils.SplitStringIntoChunks(summary, chunkSize)
	for i, chunk := range chunks {
		config.Logger.Debugf("Attempt to send chunk #%d: userId=%v, videoURL=%v", i+1, message.From.ID, videoURL)

		var msg string
		if i == 0 {
			msg = localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "telegram.result.first_message",
				TemplateData: map[string]interface{}{
					"title": videoInfo.Title,
					"text":  chunk,
				},
			})
		} else {
			msg = localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "telegram.result.message",
				TemplateData: map[string]interface{}{
					"text": chunk,
				},
			})
		}

		_, err = sendMarkdownMessage(message, msg)
		if err != nil {
			sendErrorMessage(message, localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "telegram.error.general"}))
		}
	}

	// Delete the "Processing" message
	deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, processingMsg.MessageID)
	Bot.Send(deleteMsg)
}
