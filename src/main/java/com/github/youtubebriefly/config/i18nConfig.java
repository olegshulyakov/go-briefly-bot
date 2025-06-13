package com.github.youtubebriefly.config;

import org.springframework.context.MessageSource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Primary;
import org.springframework.context.support.ReloadableResourceBundleMessageSource;
import org.springframework.context.support.YamlPropertiesLoader;

import java.util.List;
import java.util.Locale;

/**
 * Configuration class for setting up the internationalization (i18n) infrastructure.
 * This class defines a bean that loads messages from YAML files and provides a
 * consistent way to access localized messages throughout the application.
 */
@Configuration
public class i18nConfig {

    /**
     * Creates and configures a `ReloadableResourceBundleMessageSource` bean.  This bean is
     * responsible for loading message properties from YAML files, caching them for performance,
     * and providing access to localized messages.  It's marked as `@Primary` to indicate that
     * it should be the default `MessageSource` used by the application.
     *
     * @return A configured `ReloadableResourceBundleMessageSource` bean.
     */
    @Bean(name = "yamlMessageSource")
    @Primary
    public MessageSource getYamlMessageSource() {
        ReloadableResourceBundleMessageSource messageSource = new ReloadableResourceBundleMessageSource();

        // Set the base name for the message resource bundles.
        messageSource.setBasename("classpath:i18n/messages");

        // Cache messages for 1 hour (3600 seconds) to improve performance.
        messageSource.setCacheSeconds(3600);

        // Specify UTF-8 encoding to support a wide range of characters.
        messageSource.setDefaultEncoding("UTF-8");

        // Do not fall back to the system locale.  Explicitly define supported locales in the YAML files.
        messageSource.setFallbackToSystemLocale(false);

        // Set the default locale to English.
        messageSource.setDefaultLocale(Locale.ENGLISH);

        // Use a custom properties persister to handle YAML parsing.
        messageSource.setPropertiesPersister(new YamlPropertiesLoader());

        // Specify the file extensions to look for (YAML).
        messageSource.setFileExtensions(List.of(".yml", ".yaml"));

        return messageSource;
    }
}
