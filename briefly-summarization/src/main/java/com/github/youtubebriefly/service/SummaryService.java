package com.github.youtubebriefly.service;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.openai.api.common.OpenAiApiClientErrorException;
import org.springframework.retry.annotation.Retryable;
import org.springframework.stereotype.Service;

import javax.validation.constraints.NotNull;

/**
 * Service responsible for generating summaries of text using a Large Language Model (LLM).
 * It leverages an i18n service for localization and a ChatClient to interact with the LLM.
 */
@Service
public class SummaryService {

    private static final Logger logger = LoggerFactory.getLogger(SummaryService.class);

    private final i18nService i18nService;
    private final ChatClient chatClient;

    /**
     * Constructs a SummaryService instance.
     *
     * @param i18nService       The i18n service for localized messages.
     * @param chatClientBuilder A builder for creating the ChatClient instance.  This allows for
     *                          flexible configuration of the client.
     * @throws IllegalArgumentException if either dependency is null.
     */
    public SummaryService(@NotNull(message = "i18nService cannot be null") i18nService i18nService, @NotNull(message = "Builder cannot be null") ChatClient.Builder chatClientBuilder) {
        this.i18nService = i18nService;
        this.chatClient = chatClientBuilder.build();
    }

    /**
     * Generates a summary of the provided text using the LLM.
     *
     * @param text         The text to summarize.
     * @param languageCode The language code for the summary (e.g., "en", "es", "fr").
     * @return The generated summary.
     * @throws OpenAiApiClientErrorException if the LLM API returns an error after multiple retries.
     */
    @Retryable(retryFor = OpenAiApiClientErrorException.class, maxAttempts = 3)
    public String generateSummary(String text, String languageCode) {
        logger.info("Summarizing ***text*** on language: {}", languageCode);

        return chatClient.prompt()
                .user(i18nService.getMessage("llm.prompt", languageCode, text))
                .call()
                .content();
    }
}
