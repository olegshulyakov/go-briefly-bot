package com.github.youtubebriefly.service;

import lombok.RequiredArgsConstructor;
import org.springframework.context.MessageSource;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.util.Locale;

/**
 * Service for retrieving internationalized messages. This service leverages a
 * MessageSource (configured to load messages from YAML files) to provide
 * localized messages based on a given key and locale.
 */
@Service
@RequiredArgsConstructor
public class i18nService {

    private final MessageSource yamlMessageSource;

    /**
     * Retrieves an internationalized message for the given key and locale.
     *
     * @param key      The message key.
     * @param language The language code (e.g., "en", "fr", "es").
     * @param args     Optional arguments to be substituted into the message.
     * @return The localized message.  Returns the message key as a fallback if not found.
     * @throws IllegalArgumentException if key or language is null or empty.
     */
    public String getMessage(String key, String language, Object... args) {
        if (StringUtils.hasText(key)) {
            throw new IllegalArgumentException("Message key cannot be null or empty.");
        }
        if (StringUtils.hasText(language)) {
            throw new IllegalArgumentException("Language code cannot be null or empty.");
        }

        Locale locale = Locale.of(language);
        return yamlMessageSource.getMessage(key, args, locale);
    }
}
