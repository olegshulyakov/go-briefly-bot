package com.github.youtubebriefly.model;

import java.time.Instant;
import java.time.ZoneOffset;

/**
 * Represents the response containing a video transcript.  This record encapsulates
 * the type, video ID, language, transcript content, and creation timestamp.
 * Records are immutable, making them suitable for representing data that shouldn't
 * change after creation.
 */
public record VideoTranscriptResponse(
        String type,
        String videoId,
        String language,
        String transcript,
        Instant createdAt
) {
    /**
     * Constructs a {@code VideoTranscriptResponse} from a {@code VideoTranscript} object.
     * Converts the {@code VideoTranscript} object's creation timestamp to an Instant in UTC.
     *
     * @param videoTranscript The {@code VideoTranscript} object to convert.
     */
    public VideoTranscriptResponse(VideoTranscript videoTranscript) {
        this(videoTranscript.getType(), videoTranscript.getVideoId(), videoTranscript.getLanguage(), videoTranscript.getTranscript(), videoTranscript.getCreatedAt().toInstant(ZoneOffset.UTC));
    }
}
