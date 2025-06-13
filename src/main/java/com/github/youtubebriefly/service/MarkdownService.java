
package com.github.youtubebriefly.service;

import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;

/**
 * A service class responsible for converting Markdown text into HTML format.
 * This class provides a simple, custom implementation of Markdown to HTML conversion,
 * supporting common formatting elements such as headers, bold, italic, underline,
 * strikethrough, code blocks, links, and images.
 */
@Component
public class MarkdownService {

    /**
     * Converts the given Markdown string into an equivalent HTML string.
     *
     * <p>This method supports the following Markdown syntax:
     * <ul>
     *   <li>Headers: # Header</li>
     *   <li>Bold: **bold**, __bold__</li>
     *   <li>Underline: __underline__</li>
     *   <li>Italic: *italic*, _italic_</li>
     *   <li>Strikethrough: ~~strikethrough~~</li>
     *   <li>Code: `code`, ```code```, ```language\ncode```</li>
     *   <li>Links: [text](url)</li>
     *   <li>Images: ![alt](src) - converted to HTML image tag</li>
     * </ul>
     *
     * @param markdown The input Markdown string to be converted.
     * @return The resulting HTML string. If the input is null or empty, it returns the original input.
     */
    public String toHtml(String markdown) {
        if (!StringUtils.hasText(markdown)) {
            return markdown;
        }
        String html = markdown.trim();

        // Convert # Header
        html = html.replaceAll("(#+)\\s*(.+)", "$1 <strong>$2</strong>");

        // Convert **bold**, __bold__ to <strong>bold</strong>
        html = html.replaceAll("\\*\\*(.+?)\\*\\*", "<strong>$1</strong>");

        // Convert __underline__ to <ins>underline</ins>
        html = html.replaceAll("__(.+?)__", "<ins>$1</ins>");

        // Convert *italic*, _italic_ to <em>italic</em>
        html = html.replaceAll("\\*(.+?)\\*", "<em>$1</em>")
                .replaceAll("_(.+?)_", "<em>$1</em>");

        // Convert ~~strikethrough~~ to <strike>strikethrough</strike>
        html = html.replaceAll("~~(.+?)~~", "<strike>$1</strike>");

        // Convert `code`, ```code``` to <code>code</code>
        html = html.replaceAll("```([a-z]+)\\n([\\s\\S]*?)```", "<pre><code class=\"language-$1\">$2</code></pre>")
                .replaceAll("```[a-z]*\\n([\\s\\S]*?)```", "<pre><code>$1</code></pre>")
                .replaceAll("`(.+?)`", "<code>$1</code>");

        // Process images ![alt](src) - Telegram doesn't support Markdown images, convert to text
        html = html.replaceAll("!\\[([^\\]]*)\\]\\(([^\\)]+)\\)", "<img src=\"$2\" alt=\"$1\" />");

        // Process [text](url) links
        html = html.replaceAll("\\[([^\\]]+)\\]\\(([^\\)]+)\\)", "<a href=\"$2\">$1</a>");

        return html;
    }
}
