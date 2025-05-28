package com.github.youtubebriefly.config;

import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.openai.OpenAiChatModel;
import org.springframework.ai.openai.OpenAiChatOptions;
import org.springframework.ai.openai.api.OpenAiApi;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.MessageSource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import org.springframework.context.support.ReloadableResourceBundleMessageSource;
import org.springframework.context.support.YamlPropertiesLoader;

import javax.validation.constraints.NotNull;
import java.util.List;
import java.util.Locale;

@Configuration
public class AppConfig {
    @NotNull(message = "Base API URL is required, set OPENAI_BASE_URL environment variable")
    @Value("${OPENAI_BASE_URL:https://api.openai.com/v1}")
    private String openAiBaseUrl;

    @NotNull(message = "Token is required, set OPENAI_API_KEY environment variable")
    @Value("${OPENAI_API_KEY:}")
    private String openAiApiToken;

    @NotNull(message = "Model is required, set OPENAI_MODEL environment variable")
    @Value("${OPENAI_MODEL:}")
    private String openAiModel;

    @Value("${YT_DLP_PROXY:}")
    private String ytDlpProxy;

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

    @Bean
    public String getYtDlpProxy() {
        return this.ytDlpProxy;
    }

    @Bean(name = "yamlMessageSource")
    @Primary
    public MessageSource getYamlMessageSource() {
        ReloadableResourceBundleMessageSource messageSource = new ReloadableResourceBundleMessageSource();
        messageSource.setBasename("classpath:i18n/messages");
        messageSource.setCacheSeconds(3600);
        messageSource.setDefaultEncoding("UTF-8");
        messageSource.setFallbackToSystemLocale(false);
        messageSource.setDefaultLocale(Locale.ENGLISH);
        messageSource.setPropertiesPersister(new YamlPropertiesLoader());
        messageSource.setFileExtensions(List.of(".yml", ".yaml"));
        return messageSource;
    }
}
