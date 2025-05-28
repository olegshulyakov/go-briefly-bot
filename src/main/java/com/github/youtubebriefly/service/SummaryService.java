package com.github.youtubebriefly.service;

import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.openai.api.common.OpenAiApiClientErrorException;
import org.springframework.retry.annotation.Retryable;
import org.springframework.stereotype.Service;

@Service
public class SummaryService {

    private final I18nService i18nService;
    private final ChatClient chatClient;

    public SummaryService(I18nService i18nService, ChatClient.Builder chatClientBuilder) {
        this.i18nService = i18nService;
        this.chatClient = chatClientBuilder.build();
    }

    @Retryable(retryFor = OpenAiApiClientErrorException.class, maxAttempts = 3)
    public String generateSummary(String text, String languageCode) {
        return chatClient.prompt()
                .system(i18nService.getMessage("llm.system_prompt", languageCode))
                .user(text)
                .call()
                .content();
    }
}
