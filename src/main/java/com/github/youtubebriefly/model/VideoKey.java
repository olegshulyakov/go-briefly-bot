package com.github.youtubebriefly.model;

/**
 * Represents a composite key for identifying a video resource, potentially for caching or data access.
 * The key consists of the video type, the video ID, and the language.  This allows for uniquely
 * identifying different versions (language) of the same video.
 *
 * @param type The type of video (e.g., "movie", "trailer", "clip").  Must not be null or empty.
 * @param videoId The unique identifier for the video.  Must not be null or empty.
 * @param language The language of the video (e.g., "en", "es", "fr"). Must not be null or empty.
 */
public record VideoKey(String type, String videoId, String language) {
}
