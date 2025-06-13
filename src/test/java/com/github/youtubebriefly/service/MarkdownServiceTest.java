package com.github.youtubebriefly.service;

import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import static org.junit.jupiter.api.Assertions.assertEquals;

class MarkdownServiceTest {

    private final MarkdownService service = new MarkdownService();

    @ParameterizedTest
    @CsvSource({
            "This is *Sparta*,'This is <em>Sparta</em>'",

            // Basic formatting
            "'**bold**', '<strong>bold</strong>'",
            "'*italic*', '<em>italic</em>'",
            "'_italic_', '<em>italic</em>'",
            "'__underline__', '<ins>underline</ins>'",
            "'~~strikethrough~~', '<strike>strikethrough</strike>'",

            // Code blocks and inline code
            "'`code`', '<code>code</code>'",
            "'```\nmulti\nline\ncode\n```', '<pre><code>multi\nline\ncode\n</code></pre>'",
            "'```java\nSystem.out.println();\n```', '<pre><code class=\"language-java\">System.out.println();\n</code></pre>'",

            // Links and images
            "'[text](https://example.com)', '<a href=\"https://example.com\">text</a>'",
            "'![alt](image.png)', '<img src=\"image.png\" alt=\"alt\" />'",

            // Headers
            "'# Header 1', '# <strong>Header 1</strong>'",
            "'## Header 2', '## <strong>Header 2</strong>'",
            "'### Header 3', '### <strong>Header 3</strong>'",

            // Edge cases
            "'', ''",  // empty string
            "'plain text', 'plain text'",  // plain text
            "'multiple    spaces', 'multiple    spaces'"  // preserve whitespace
    })
    void toHtml(String markdown, String html) {
        assertEquals(html, service.toHtml(markdown));
    }
}