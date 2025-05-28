package com.github.youtubebriefly.service;

import org.springframework.ai.chat.client.ChatClient;
import org.springframework.retry.annotation.Retryable;
import org.springframework.stereotype.Service;

@Service
public class SummaryService {

    private final ChatClient chatClient;

    public SummaryService(ChatClient.Builder chatClientBuilder) {
        this.chatClient = chatClientBuilder.build();
    }

    @Retryable(retryFor = RuntimeException.class, maxAttempts = 3)
    public String generateSummary(String text) {
        return chatClient.prompt()
                .system("")
                .user(text)
                .call()
                .content();
    }
}
