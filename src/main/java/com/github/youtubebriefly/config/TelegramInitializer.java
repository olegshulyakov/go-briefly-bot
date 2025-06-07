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

@Component
@RequiredArgsConstructor
public class TelegramInitializer {
    private static final Logger logger = LoggerFactory.getLogger(TelegramInitializer.class);

    private final TelegramBotController telegramBot;

    @PostConstruct
    public void afterPropertiesSet() {
        try {
            TelegramBotsApi telegramBotsApi = new TelegramBotsApi(DefaultBotSession.class);
            telegramBotsApi.registerBot(telegramBot);
        } catch (TelegramApiException e) {
            logger.error("Cannot register Telegram Bot", e);
        }
    }
}
