package com.github.youtubebriefly.service;

import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import static org.junit.jupiter.api.Assertions.assertEquals;

class MarkdownServiceTest {

    private final MarkdownService service = new MarkdownService();

    @ParameterizedTest
    @CsvSource({
            "This is *Sparta*,'<p>This is <em>Sparta</em></p>\n'"
    })
    void toHtml(String markdown, String html) {
        assertEquals(html, service.toHtml(markdown));
    }
}