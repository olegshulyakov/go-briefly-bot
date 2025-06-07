package com.github.youtubebriefly.config;

import lombok.Getter;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;

/**
 * Configuration class for Telegram bot properties.
 */
@Configuration
@Getter
public class TelegramConfig {

    /**
     * The Telegram bot token.  This is a secret key used to interact with the Telegram Bot API.
     */
    @Value("${TELEGRAM_BOT_TOKEN:}")
    private String botToken;

    /**
     * The Telegram bot username.  Used to identify the bot on Telegram.
     */
    @Value("${TELEGRAM_BOT_USERNAME:}")
    private String botUsername;
}
