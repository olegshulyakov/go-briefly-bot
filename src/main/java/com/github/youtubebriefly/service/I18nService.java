package com.github.youtubebriefly.service;

import lombok.RequiredArgsConstructor;
import org.springframework.context.MessageSource;
import org.springframework.stereotype.Service;

import java.util.Locale;

@Service
@RequiredArgsConstructor
public class I18nService {

    private final MessageSource yamlMessageSource;

    public String getMessage(String key, String language, Object... args) {
        Locale locale = Locale.of(language);
        return yamlMessageSource.getMessage(key, args, locale);
    }
}
