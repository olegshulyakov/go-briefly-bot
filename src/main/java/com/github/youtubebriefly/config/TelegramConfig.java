package com.github.youtubebriefly.config;

import lombok.Getter;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;

@Configuration
@Getter
public class TelegramConfig {
    @Value("${TELEGRAM_BOT_TOKEN:}")
    private String botToken;

    @Value("${TELEGRAM_BOT_USERNAME:}")
    private String botUsername;
}
