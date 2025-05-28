package com.github.youtubebriefly.model;

import java.time.Instant;
import java.time.ZoneOffset;

public record VideoTranscriptResponse(
        String type,
        String videoId,
        String languageCode,
        String transcript,
        Instant createdAt
) {
    public VideoTranscriptResponse(VideoTranscript videoTranscript) {
        this(videoTranscript.getType(), videoTranscript.getVideoId(), videoTranscript.getLanguageCode(), videoTranscript.getTranscript(), videoTranscript.getCreatedAt().toInstant(ZoneOffset.UTC));
    }
}
