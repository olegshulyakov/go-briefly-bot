package com.github.youtubebriefly.config;

import com.github.youtubebriefly.controller.TelegramBotController;
import jakarta.annotation.PostConstruct;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import org.telegram.telegrambots.meta.TelegramBotsApi;
import org.telegram.telegrambots.meta.exceptions.TelegramApiException;
import org.telegram.telegrambots.updatesreceivers.DefaultBotSession;

/**
 * Component responsible for initializing and registering the Telegram bot with the Telegram Bots API.
 * It uses the {@link TelegramBotsApi} to register the bot instance.  Handles potential
 * {@link TelegramApiException} during registration and logs any errors encountered.
 */
@Component
@RequiredArgsConstructor
public class TelegramInitializer {

    private static final Logger logger = LoggerFactory.getLogger(TelegramInitializer.class);

    private final TelegramBotController telegramBot;

    /**
     * Initializes and registers the Telegram bot with the Telegram Bots API after the component
     * has been fully constructed and all dependencies have been injected.
     * This method runs after the Spring context is initialized.
     */
    @PostConstruct
    public void initialize() {
        try {
            TelegramBotsApi telegramBotsApi = new TelegramBotsApi(DefaultBotSession.class);
            telegramBotsApi.registerBot(telegramBot);
            logger.info("Telegram bot registered successfully.");
        } catch (TelegramApiException e) {
            logger.error("Failed to register Telegram Bot.", e);
        }
    }
}
