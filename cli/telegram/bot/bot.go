// Package bot provides functionality to handle Telegram bot interactions,
// including message processing, rate limiting, and YouTube video summarization.
//
// The package offers:
// - Telegram bot initialization and configuration
// - Message handling with rate limiting
// - YouTube video URL extraction and processing
// - Video transcript fetching and summarization
// - Localized responses based on user language
// - Message formatting and chunking for Telegram's message length limits
//
// Key Features:
// - Rate limiting to prevent abuse (30-second cooldown per user)
// - Support for both plain text and formatted (Markdown/HTML) messages
// - Automatic splitting of long messages to comply with Telegram's limits
// - Progress updates during video processing
// - Error handling with localized user feedback
package bot

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	tg_md2html "github.com/PaulSonOfLars/gotg_md2html"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/olegshulyakov/go-briefly-bot/lib"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/youtube"
	"github.com/olegshulyakov/go-briefly-bot/lib/transformers/summarization"
)

const (
	// maxLength defines the maximum length for a single Telegram message.
	maxLength = 4000
)

var (
	// Bot is the global Telegram bot instance.
	Bot *tgbotapi.BotAPI

	// userLastRequest tracks the last request time for each user to enforce rate limiting.
	userLastRequest = make(map[int64]time.Time)

	// userMutex is a mutex to protect concurrent access to the userLastRequest map.
	userMutex = sync.Mutex{}
)

// New initializes a new Telegram bot instance with the provided API token.
// It sets the global Bot variable and returns the initialized bot instance.
func New(token string) (*tgbotapi.BotAPI, error) {
	var err error
	Bot, err = tgbotapi.NewBotAPI(token)
	return Bot, err
}

// Handle processes incoming Telegram updates, routing them to appropriate handlers.
// It performs initial validation including:
// - Nil message checks
// - Bot user detection
// - Rate limiting enforcement
// - Command handling (currently only "/start")
// - Regular message processing.
func Handle(update tgbotapi.Update) {
	message := update.Message
	if message == nil {
		return
	}

	// Check if the user is bot
	if message.From.IsBot {
		slog.Warn(
			"Got message from bot",
			"userID",
			message.From.ID,
			"user",
			message.From,
			"bot",
			message.From.IsBot,
			"language",
			message.From.LanguageCode,
		)
		return
	}

	// Check if the user is rate-limited
	if isUserRateLimited(message.From.ID) {
		slog.Warn(
			"Rate Limit exceeded",
			"userId",
			message.From.ID,
			"user",
			message.From,
			"language",
			message.From.LanguageCode,
		)
		sendQuite(message, lib.MustLocalize(message.From.LanguageCode, "telegram.error.rate_limited"))
		return
	}

	// Handle commands
	if message.IsCommand() {
		if message.Command() == "start" {
			_, _ = send(message, lib.MustLocalize(message.From.LanguageCode, "telegram.welcome.message"))
		}
		return
	}

	handle(message)
}

// sendRetry attempts to send a message with retry logic (max 3 attempts).
// Returns the sent message or error if all attempts fail.
func sendRetry(msg tgbotapi.Chattable) (tgbotapi.Message, error) {
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

// send creates and sends a plain text reply message to the user.
// Returns the sent message or error.
func send(userMessage *tgbotapi.Message, text string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(userMessage.Chat.ID, text)
	msg.ReplyToMessageID = userMessage.MessageID
	return sendRetry(msg)
}

// sendFormatted creates and sends a formatted (HTML) reply message to the user.
// Returns the sent message or error.
func sendFormatted(userMessage *tgbotapi.Message, text string) (tgbotapi.Message, error) {
	escapedMessageText := tg_md2html.MD2HTMLV2(text)

	msg := tgbotapi.NewMessage(userMessage.Chat.ID, escapedMessageText)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true
	return sendRetry(msg)
}

// edit modifies an existing message with new text.
// Returns the edited message or error.
func edit(userMessage *tgbotapi.Message, messageToEdit tgbotapi.Message, text string) (tgbotapi.Message, error) {
	editMsg := tgbotapi.NewEditMessageText(userMessage.Chat.ID, messageToEdit.MessageID, text)
	return sendRetry(editMsg)
}

// sendQuite sends a message without returning any response or error (fire-and-forget).
func sendQuite(userMessage *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(userMessage.Chat.ID, text)
	msg.ReplyToMessageID = userMessage.MessageID
	_, _ = Bot.Send(msg)
}

// deleteQuite deletes a message without returning any response or error (fire-and-forget).
func deleteQuite(userMessage *tgbotapi.Message, message tgbotapi.Message) {
	deleteMsg := tgbotapi.NewDeleteMessage(userMessage.Chat.ID, message.MessageID)
	_, _ = Bot.Send(deleteMsg)
}

// isUserRateLimited checks if a user has made a request within the rate limit window (30 seconds).
// Updates the last request time if the user is not rate-limited.
// Returns true if the user should be rate-limited.
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

// handle processes non-command messages containing YouTube URLs.
// It performs the following steps:
// 1. Extracts YouTube URLs from the message
// 2. Sends progress updates to the user
// 3. Fetches video transcript
// 4. Generates a summary
// 5. Sends the results back to the user
// 6. Cleans up progress messages.
func handle(message *tgbotapi.Message) {
	text := message.Text
	slog.Debug("Telegram: Request", "userId", message.From.ID, "user", message.From, "language", message.From.LanguageCode, "text", text)

	// Check if the message contains a YouTube link and extract URL
	videoURLs := youtube.ExtractURLs(text)
	if len(videoURLs) == 0 {
		slog.Error("Got invalid processing message", "userId", message.From.ID, "user", message.From, "text", text)
		sendQuite(message, lib.MustLocalize(message.From.LanguageCode, "telegram.error.no_url_found"))
		return
	}

	// Notify user about process start
	slog.Info("Processing YouTube video", "urls", videoURLs)
	processingMsg, err := send(message, lib.MustLocalize(message.From.LanguageCode, "telegram.progress.processing"))
	if err != nil {
		slog.Error("Failed to send processing message: %v", "error", err)
		return
	}

	// Check if there are multiple URLs
	if len(videoURLs) > 1 {
		processingMsg, _ = edit(message, processingMsg, lib.MustLocalize(message.From.LanguageCode, "telegram.error.multiple_urls"))
	}
	videoURL := videoURLs[0]

	// Fetch video info
	processingMsg, err = edit(message, processingMsg, lib.MustLocalize(message.From.LanguageCode, "telegram.progress.fetching_info"))
	if err != nil {
		slog.Error("Failed to update progress message: %v", "error", err)
		return
	}

	loader, err := loaders.VideoLoader(videoURL)
	if err == nil {
		err = loader.Load()
	}
	if err != nil {
		slog.Error("Failed to get transcript", "userId", message.From.ID, "videoURL", videoURL, "error", err)
		_, _ = edit(message, processingMsg, lib.MustLocalize(message.From.LanguageCode, "telegram.error.transcript_failed"))
		return
	}
	videoTranscript := loader.Transcript()

	// Summarize transcript
	processingMsg, err = edit(message, processingMsg, lib.MustLocalize(message.From.LanguageCode, "telegram.progress.summarizing"))
	if err != nil {
		slog.Error("Failed to update progress message", "userId", message.From.ID, "error", err)
		return
	}

	summary, err := summarization.SummarizeText(videoTranscript.Transcript, message.From.LanguageCode)
	if err != nil {
		slog.Error("Failed to summarize transcript", "userId", message.From.ID, "videoURL", videoURL, "error", err)
		_, _ = edit(message, processingMsg, lib.MustLocalize(message.From.LanguageCode, "telegram.error.summary_failed"))
		return
	}

	// Send response to user
	response := fmt.Sprintf(
		"%s\n%s",
		lib.MustLocalizeTemplate(
			message.From.LanguageCode,
			"telegram.response.title",
			map[string]string{"title": videoTranscript.Title},
		),
		summary,
	)
	chunks := lib.ToChunks(response, maxLength)
	for i, chunk := range chunks {
		slog.Debug("Attempt to send chunk", "chunk", i+1, "userId", message.From.ID, "videoURL", videoURL)
		_, err = sendFormatted(message, chunk)
		if err != nil {
			slog.Error("Failed to send chunk", "chunk", i+1, "userId", message.From.ID, "videoURL", videoURL, "error", err)
		}
	}

	// Delete the "Processing" message at the end
	deleteQuite(message, processingMsg)
}
