// Package bot provides functionality to handle Telegram bot interactions,
// including message processing, rate limiting, and YouTube video summarization.
package bot

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	tg_md2html "github.com/PaulSonOfLars/gotg_md2html"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/olegshulyakov/go-briefly-bot/lib"
	"github.com/olegshulyakov/go-briefly-bot/lib/summarization"
	"github.com/olegshulyakov/go-briefly-bot/lib/transcript"
	"github.com/olegshulyakov/go-briefly-bot/lib/transcript/youtube"
)

const telegramMessageMaxLength = 4000

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
		slog.Error("Failed to create bot: %v", "error", err)
	}

	// Bot.Debug = true
	slog.Info("Bot started successfully")

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
	slog.Info("Shutting down bot...")
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
		slog.Warn("Failed to send message", "attempt", i+1, "error", err)
	}

	slog.Error(fmt.Sprintf("Failed to send message after %d attempts", maxRetries), "error", err)
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
	_, _ = Bot.Send(msg)
}

// isUserRateLimited checks if the user has made a request within the last 30 seconds.
//
// Parameters:
//   - userId: The ID of the user to check.
//
// Returns:
//   - bool: True if the user is rate-limited, false otherwise.
func isUserRateLimited(userID int64) bool {
	userMutex.Lock()
	defer userMutex.Unlock()

	lastRequest, exists := userLastRequest[userID]
	if exists && time.Since(lastRequest) < 30*time.Second {
		return true // User is rate-limited
	}

	userLastRequest[userID] = time.Now() // Update the last request time
	return false
}

// handleTelegramMessage processes incoming messages from users, including commands and YouTube links.
//
// Parameters:
//   - message: The incoming message to process.
func handleTelegramMessage(message *tgbotapi.Message) {
	// Check if the user is bot
	if message.From.IsBot {
		slog.Warn("Got message from bot", "userId", message.From.ID, "user", message.From, "bot", message.From.IsBot, "language", message.From.LanguageCode)
		return
	}

	// Check if the user is rate-limited
	if isUserRateLimited(message.From.ID) {
		slog.Warn("Rate Limit exceeded", "userId", message.From.ID, "user", message.From, "language", message.From.LanguageCode)
		sendErrorMessage(message, lib.MustLocalize(message.From.LanguageCode, "telegram.error.rate_limited"))
		return
	}

	text := message.Text
	slog.Debug("Telegram: Request", "userId", message.From.ID, "user", message.From, "language", message.From.LanguageCode, "text", text)

	if message.IsCommand() {
		if message.Command() == "start" {
			_, _ = sendMessage(message, lib.MustLocalize(message.From.LanguageCode, "telegram.welcome.message"))
		}
		return
	}

	// Check if the message contains a YouTube link and extract URL
	videoURLs, err := youtube.ExtractURLs(text)
	if err != nil {
		slog.Error("Got invalid processing message", "userId", message.From.ID, "user", message.From, "text", text)
		sendErrorMessage(message, lib.MustLocalize(message.From.LanguageCode, "telegram.error.no_url_found"))
		return
	}

	// Notify user about process start
	slog.Info("Processing YouTube video", "urls", videoURLs)
	processingMsg, err := sendMessage(message, lib.MustLocalize(message.From.LanguageCode, "telegram.progress.processing"))
	if err != nil {
		slog.Error("Failed to send processing message: %v", "error", err)
		return
	}

	// Check if there are multiple URLs
	if len(videoURLs) > 1 {
		processingMsg, _ = editMessage(message, processingMsg, lib.MustLocalize(message.From.LanguageCode, "telegram.error.multiple_urls"))
	}
	videoURL := videoURLs[0]

	// Fetch video info
	processingMsg, err = editMessage(message, processingMsg, lib.MustLocalize(message.From.LanguageCode, "telegram.progress.fetching_info"))
	if err != nil {
		slog.Error("Failed to update progress message: %v", "error", err)
		return
	}

	videoTranscript, err := transcript.Transcript(videoURL)
	if err != nil {
		slog.Error("Failed to get transcript", "userId", message.From.ID, "videoURL", videoURL, "error", err)
		_, _ = editMessage(message, processingMsg, lib.MustLocalize(message.From.LanguageCode, "telegram.error.transcript_failed"))
		return
	}

	// Summarize transcript
	processingMsg, err = editMessage(message, processingMsg, processingMsg.Text+"\n"+lib.MustLocalize(message.From.LanguageCode, "telegram.progress.summarizing"))
	if err != nil {
		slog.Error("Failed to update progress message", "userId", message.From.ID, "error", err)
		return
	}

	summary, err := summarization.SummarizeText(videoTranscript.Transcript, message.From.LanguageCode)
	if err != nil {
		slog.Error("Failed to summarize transcript", "userId", message.From.ID, "videoURL", videoURL, "error", err)
		_, _ = editMessage(message, processingMsg, lib.MustLocalize(message.From.LanguageCode, "telegram.error.summary_failed"))
		return
	}

	// Send summary to user
	chunkSize := telegramMessageMaxLength - len(videoTranscript.Title)
	chunks := lib.SplitStringIntoChunks(summary, chunkSize)
	for i, chunk := range chunks {
		slog.Debug("Attempt to send chunk", "chunk", i+1, "userId", message.From.ID, "videoURL", videoURL)

		var msg string
		if i == 0 {
			msg = lib.MustLocalizeTemplate(
				message.From.LanguageCode,
				"telegram.result.first_message",
				map[string]string{
					"title": videoTranscript.Title,
					"text":  chunk,
				},
			)
		} else {
			msg = lib.MustLocalizeTemplate(
				message.From.LanguageCode,
				"telegram.result.message",
				map[string]string{
					"text": chunk,
				},
			)
		}

		_, err = sendMarkdownMessage(message, msg)
		if err != nil {
			sendErrorMessage(message, lib.MustLocalize(message.From.LanguageCode, "telegram.error.general"))
		}
	}

	// Delete the "Processing" message
	deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, processingMsg.MessageID)
	_, _ = Bot.Send(deleteMsg)
}
