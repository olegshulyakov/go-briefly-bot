package com.github.youtubebriefly.model;

/**
 * Represents a summary response, containing the summary text and the language in which it was generated.
 */
public record SummaryResponse(String summary, String language) {
}
