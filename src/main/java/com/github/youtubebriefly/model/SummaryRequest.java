package com.github.youtubebriefly.model;

import javax.validation.constraints.NotNull;

public record SummaryRequest(@NotNull(message = "Text is required") String text, String language) {
}
