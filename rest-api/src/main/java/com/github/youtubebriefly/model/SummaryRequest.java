package com.github.youtubebriefly.model;

import jakarta.validation.constraints.NotNull;

/**
 * Represents a request to summarize a given text.  This record encapsulates
 * the text to be summarized and the desired language code for the summary.
 */
public record SummaryRequest(@NotNull(message = "Text is required") String text, String languageCode) {
}
