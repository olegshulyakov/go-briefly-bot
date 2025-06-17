package com.github.youtubebriefly.config;

import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.openai.OpenAiChatModel;
import org.springframework.ai.openai.OpenAiChatOptions;
import org.springframework.ai.openai.api.OpenAiApi;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import javax.validation.constraints.NotNull;

/**
 * Configuration class for OpenAI integration.  This class manages the configuration
 * of the OpenAI API client, including the base URL, API token, and default model.
 * It uses Spring's `@Configuration` annotation to mark it as a source of bean definitions.
 */
@Configuration
public class OpenAiConfig {

    /**
     * The base URL for the OpenAI API.  This is read from the `OPENAI_BASE_URL` environment variable.
     * Defaults to OpenAI API if the environment variable is not set.
     */
    @NotNull(message = "Base API URL is required, set OPENAI_BASE_URL environment variable")
    @Value("${OPENAI_BASE_URL:https://api.openai.com/v1}")
    private String openAiBaseUrl;

    /**
     * The OpenAI API token.  This is read from the `OPENAI_API_KEY` environment variable.
     * This token is required for authentication.
     */
    @NotNull(message = "Token is required, set OPENAI_API_KEY environment variable")
    @Value("${OPENAI_API_KEY:}")
    private String openAiApiToken;

    /**
     * The default OpenAI model to use.  This is read from the `OPENAI_MODEL` environment variable.
     * A model must be specified.
     */
    @NotNull(message = "Model is required, set OPENAI_MODEL environment variable")
    @Value("${OPENAI_MODEL:}")
    private String openAiModel;

    /**
     * Creates and configures a `ChatClient.Builder`.  This builder is used to create
     * instances of the `ChatClient` with the configured OpenAI API.
     *
     * @return A `ChatClient.Builder` instance configured with the OpenAI API settings.
     */
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
