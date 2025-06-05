package com.github.youtubebriefly.config;

import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.openai.OpenAiChatModel;
import org.springframework.ai.openai.OpenAiChatOptions;
import org.springframework.ai.openai.api.OpenAiApi;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import javax.validation.constraints.NotNull;

@Configuration
public class OpenAiConfig {
    @NotNull(message = "Base API URL is required, set OPENAI_BASE_URL environment variable")
    @Value("${OPENAI_BASE_URL:https://api.openai.com/v1}")
    private String openAiBaseUrl;

    @NotNull(message = "Token is required, set OPENAI_API_KEY environment variable")
    @Value("${OPENAI_API_KEY:}")
    private String openAiApiToken;

    @NotNull(message = "Model is required, set OPENAI_MODEL environment variable")
    @Value("${OPENAI_MODEL:}")
    private String openAiModel;

    @Bean
    public ChatClient.Builder getChatClientBuilder() {
        return ChatClient.builder(
                OpenAiChatModel.builder().openAiApi(
                        OpenAiApi.builder()
                                .baseUrl(openAiBaseUrl)
                                .apiKey(openAiApiToken)
                                .build()
                ).defaultOptions(
                        OpenAiChatOptions.builder()
                                .model(openAiModel)
                                .build()
                ).build()
        );
    }
}
