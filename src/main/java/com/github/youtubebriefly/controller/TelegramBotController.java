package com.github.youtubebriefly.controller;

import com.github.youtubebriefly.config.TelegramConfig;
import com.github.youtubebriefly.dao.UserRequestRepository;
import com.github.youtubebriefly.exception.YouTubeException;
import com.github.youtubebriefly.model.UserRequest;
import com.github.youtubebriefly.model.VideoSummary;
import com.github.youtubebriefly.service.i18nService;
import com.github.youtubebriefly.service.YouTubeService;
import com.github.youtubebriefly.service.YouTubeSummaryService;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.ai.openai.api.common.OpenAiApiClientErrorException;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.bots.TelegramLongPollingBot;
import org.telegram.telegrambots.meta.api.methods.ParseMode;
import org.telegram.telegrambots.meta.api.methods.send.SendMessage;
import org.telegram.telegrambots.meta.api.methods.updatingmessages.DeleteMessage;
import org.telegram.telegrambots.meta.api.methods.updatingmessages.EditMessageText;
import org.telegram.telegrambots.meta.api.objects.Message;
import org.telegram.telegrambots.meta.api.objects.Update;
import org.telegram.telegrambots.meta.api.objects.User;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * A Telegram bot controller that handles incoming updates and interacts with the Telegram API.
 * It extracts YouTube video URLs from messages, summarizes the videos using a separate service,
 * and sends the summaries back to the user in chunks.  It includes rate limiting and error handling.
 */
@Component
@RequiredArgsConstructor
public class TelegramBotController extends TelegramLongPollingBot {
    private static final Logger logger = LoggerFactory.getLogger(TelegramBotController.class);
    private static final int MAX_RETRIES = 3;
    private static final int RATE_LIMIT_SECONDS = 30;
    private static final int MESSAGE_CHUNK_SIZE = 4000;

    /**
     * Stores the last request timestamp for each user to enforce rate limiting.
     * Key: User ID, Value: Last request timestamp.
     */
    private final Map<Long, LocalDateTime> userLastRequest = new HashMap<>();
    private final TelegramConfig config;
    private final UserRequestRepository userRequestRepository;
    private final i18nService i18nService;
    private final YouTubeSummaryService youTubeSummaryService;

    /**
     * {@inheritDoc}
     */
    @Override
    public String getBotUsername() {
        return config.getBotUsername();
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public String getBotToken() {
        return config.getBotToken();
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void onUpdateReceived(Update update) {
        if (!update.hasMessage()) {
            return;
        }

        Message message = update.getMessage();
        User user = message.getFrom();

        // Skip messages from bots
        if (user.getIsBot()) {
            logger.warn("Got message from bot: userId={}, user='{}', language='{}'", user.getId(), user, user.getLanguageCode());
            return;
        }

        // Check rate limiting
        if (isUserRateLimited(user.getId())) {
            logger.warn("Rate Limit exceeded: userId={}, user='{}', language='{}'", user.getId(), user, user.getLanguageCode());
            sendMessageByKey(message, "telegram.error.rate_limited");
            return;
        }

        logger.debug("Request: userId={}, user='{}', language='{}', text={}", user.getId(), user, user.getLanguageCode(), message.getText());

        if (message.isCommand()) {
            handleCommand(message);
            return;
        }

        UserRequest userRequest = userRequestRepository.save(new UserRequest(null, user.getId(), user.getLanguageCode(), message.getText(), LocalDateTime.now(), null));

        // Process message
        handleMessage(message);

        userRequest.setProceedAt(LocalDateTime.now());
        userRequestRepository.save(userRequest);
    }

    /**
     * Handles incoming command messages.
     * Currently supports the "start" command.
     *
     * @param message The incoming message.
     */
    private void handleCommand(Message message) {
        String command = message.getText().split(" ")[0].substring(1); // Remove '/'
        if (command.equals("start")) {
            sendMessageByKey(message, "telegram.welcome.message");
        } else {
            sendMessageByKey(message, "telegram.error.unknown_command");
        }
    }

    /**
     * Handles incoming regular messages (non-command messages).
     * Extracts YouTube video URLs, summarizes the videos, and sends the summaries back to the user.
     *
     * @param message The incoming message.
     */
    private void handleMessage(Message message) {
        String text = message.getText();
        User user = message.getFrom();
        String language = user.getLanguageCode();

        List<String> videoUrls;
        try {
            videoUrls = YouTubeService.extractAllUrls(text);
        } catch (Exception e) {
            logger.error("Got invalid processing message: userId={}, user='{}', text={}", user.getId(), user, text);
            sendMessageByKey(message, "telegram.error.no_url_found");
            return;
        }

        if (videoUrls.isEmpty()) {
            sendMessageByKey(message, "telegram.error.no_url_found");
            return;
        }

        logger.info("Processing YouTube video: {}", videoUrls);
        Message processingMsg = sendMessageByKey(message, "telegram.progress.processing");
        if (processingMsg == null) {
            logger.error("Error sending processing message: userId={}", user.getId());
            return;
        }

        if (videoUrls.size() > 1) {
            editMessage(processingMsg, i18nService.getMessage("telegram.error.multiple_urls", language));
        }

        String videoUrl = videoUrls.getFirst();
        VideoSummary summary;

        // Get summary
        try {
            summary = youTubeSummaryService.getSummary(videoUrl, language);
        } catch (YouTubeException e) {
            logger.error("Error fetch YouTube video info: userId={}, videoURL={}", user.getId(), videoUrl, e);
            editMessage(processingMsg, i18nService.getMessage("telegram.error.info_failed", language));
            return;
        } catch (OpenAiApiClientErrorException e) {
            logger.error("Error summarize YouTube video: userId={}, videoURL={}", user.getId(), videoUrl, e);
            editMessage(processingMsg, i18nService.getMessage("telegram.error.summary_failed", language));
            return;
        }

        try {
            // Send summary in chunks
            sendSummary(message, summary.getTitle(), summary.getSummary(), language);
            // Delete processing message
            deleteMessage(processingMsg);
        } catch (TelegramApiException e) {
            logger.error("Error processing YouTube video: userId={}, videoURL={}", user.getId(), videoUrl, e);
            editMessage(processingMsg, i18nService.getMessage("telegram.error.general", language));
        }
    }

    /**
     * Sends a summary message to the user, splitting it into chunks if necessary.
     *
     * @param originalMessage The original message from the user.
     * @param title           The title of the summary.
     * @param summary         The summary text.
     * @param language        The user's language code.
     */
    private void sendSummary(Message originalMessage, String title, String summary, String language) throws TelegramApiException {
        int chunkSize = MESSAGE_CHUNK_SIZE - title.length();
        List<String> chunks = this.splitTextIntoChunks(summary, chunkSize);

        for (int i = 0; i < chunks.size(); i++) {
            String chunk = chunks.get(i);
            String messageText;

            if (i == 0) {
                messageText = i18nService.getMessage("telegram.result.first_message", language, title, chunk);
            } else {
                messageText = i18nService.getMessage("telegram.result.message", language, chunk);
            }

            sendMarkdownMessage(originalMessage, messageText);
        }
    }

    /**
     * Checks if a user has exceeded the rate limit.
     *
     * @param userId The user's ID.
     * @return True if the user is rate-limited, false otherwise.
     */
    private boolean isUserRateLimited(Long userId) {
        LocalDateTime lastRequest = userLastRequest.get(userId);
        if (lastRequest != null && lastRequest.plusSeconds(RATE_LIMIT_SECONDS).isAfter(LocalDateTime.now())) {
            return true;
        }
        userLastRequest.put(userId, LocalDateTime.now());
        return false;
    }

    /**
     * Edits an existing message.
     *
     * @param messageToEdit The message to edit.
     * @param newText       The new message text.
     */
    private void editMessage(Message messageToEdit, String newText) {
        EditMessageText editMessage = new EditMessageText();
        editMessage.setChatId(messageToEdit.getChatId().toString());
        editMessage.setMessageId(messageToEdit.getMessageId());
        editMessage.setText(newText);

        try {
            execute(editMessage);
        } catch (TelegramApiException e) {
            logger.error("Failed to edit message", e);
        }
    }

    /**
     * Deletes an existing message.
     *
     * @param messageToDelete The message to delete.
     */
    private void deleteMessage(Message messageToDelete) {
        DeleteMessage deleteMessage = new DeleteMessage();
        deleteMessage.setChatId(messageToDelete.getChatId().toString());
        deleteMessage.setMessageId(messageToDelete.getMessageId());

        try {
            execute(deleteMessage);
        } catch (TelegramApiException e) {
            logger.error("Failed to delete message", e);
        }
    }

    /**
     * Sends a message with Markdown formatting.
     *
     * @param originalMessage The original message from the user.
     * @param text            The message text.
     */
    private void sendMarkdownMessage(Message originalMessage, String text) throws TelegramApiException {
        SendMessage message = new SendMessage();
        message.setChatId(originalMessage.getChatId().toString());
        message.setText(text);
        message.setParseMode(ParseMode.MARKDOWN);
        message.setDisableWebPagePreview(true);

        sendWithRetry(message);
    }

    /**
     * Sends a message by key of the message in the I18nService to the user.
     *
     * @param originalMessage The original message from the user.
     * @param messageKey The key of the message in the I18nService.
     * @return The sent message object, or null if sending failed.
     */
    private Message sendMessageByKey(Message originalMessage, String messageKey) {
        String text = i18nService.getMessage(messageKey, originalMessage.getFrom().getLanguageCode());
        SendMessage message = new SendMessage();
        message.setChatId(originalMessage.getChatId().toString());
        message.setText(text);
        message.setReplyToMessageId(originalMessage.getMessageId());

        try {
            return sendWithRetry(message);
        } catch (TelegramApiException e) {
            logger.error("Failed to send error message", e);
        }
        return null;
    }

    /**
     * Sends a message with retries.
     *
     * @param message The message to send.
     * @return The sent message object.
     * @throws TelegramApiException if the message could not be sent after multiple retries.
     */
    private Message sendWithRetry(SendMessage message) throws TelegramApiException {
        TelegramApiException lastException = null;

        for (int i = 0; i < MAX_RETRIES; i++) {
            try {
                return execute(message);
            } catch (TelegramApiException e) {
                lastException = e;
                logger.warn("Attempt {}: Failed to send message", i + 1, e);
                try {
                    Thread.sleep(1000 * (i + 1)); // Exponential backoff
                } catch (InterruptedException ie) {
                    Thread.currentThread().interrupt();
                }
            }
        }

        logger.error("Failed to send message after {} attempts", MAX_RETRIES, lastException);
        throw new TelegramApiException("Failed to send message after retries", lastException);
    }

    /**
     * Splits a string into chunks of a maximum length, ensuring that each chunk is
     * valid UTF-8 and doesn't break runes (characters).
     *
     * @param text      The string to split.
     * @param chunkSize The maximum size of each chunk (in runes/characters). Must be a positive integer.
     * @return A list of strings, where each string is a chunk of the original string.
     * @throws IllegalArgumentException if chunkSize is not positive.
     */
    private List<String> splitTextIntoChunks(String text, int chunkSize) {
        if (chunkSize <= 0) {
            throw new IllegalArgumentException("Chunk size must be a positive integer.");
        }

        int totalLength = text.length();
        List<String> chunks = new ArrayList<>(Math.divideExact(totalLength, chunkSize));
        for (int i = 0; i < totalLength; i += chunkSize) {
            chunks.add(text.substring(i, Math.min(i + chunkSize, totalLength)));
        }
        return chunks;
    }
}
