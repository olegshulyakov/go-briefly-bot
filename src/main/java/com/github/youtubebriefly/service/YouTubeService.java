package com.github.youtubebriefly.service;

import com.github.youtubebriefly.exception.YouTubeException;
import com.github.youtubebriefly.model.VideoInfo;
import com.github.youtubebriefly.model.VideoTranscript;

import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * YouTube Service
 * Handles video info retrieval and transcript downloading.
 */
public interface YouTubeService {
    /**
     * <p>
     * Constant string representing the regular expression pattern used to match YouTube video URLs.
     * This pattern is designed to handle various YouTube URL formats, including:
     * - Full URLs with "https://" or "http://"
     * - URLs with "www."
     * - Shortened URLs like "youtu.be/"
     * - URLs with the standard "youtube.com/watch?v="
     * </p>
     *
     * <p>
     * The pattern captures the YouTube video ID in group 1.
     * </p>
     */
    String YOUTUBE_URL_PATTERN = "^(?:https?://)?(?:www\\.)?(?:youtube\\.com/watch\\?v=|youtu\\.be/)([a-zA-Z0-9_-]{11})$";

    /**
     * Checks if the given URL is a valid YouTube video URL.
     *
     * @param url The URL to validate.  Must not be null.
     * @return {@code true} if the URL is a valid YouTube video URL, {@code false} otherwise.
     * @throws IllegalArgumentException if the input URL is null. This is more robust than just returning false,
     *                                  as it explicitly signals an invalid input condition.  This aligns with defensive programming.
     */
    default boolean isValidUrl(String url) {
        // Null check
        if (url == null) {
            return false;
        }

        // Compile the pattern into a regex object. This is done once per URL validation
        // for efficiency (the pattern itself doesn't change).
        Pattern pattern = Pattern.compile(YOUTUBE_URL_PATTERN);
        Matcher matcher = pattern.matcher(url);
        return matcher.matches();
    }

    /**
     * Parses the YouTube video ID from the given URL.
     *
     * @param url The URL to parse.  Must not be null.
     * @return An {@link String} containing the video ID if the URL is a valid YouTube video URL;
     * otherwise, an empty {@link Optional}. This avoids NullPointerExceptions and provides a cleaner way
     * to represent the possibility of a missing ID.
     * @throws YouTubeException if the input URL is null or ID not found.
     */
    default String getVideoId(String url) throws YouTubeException {
        // Null check
        if (url == null) {
            throw new YouTubeException("URL isd empty");
        }

        Pattern pattern = Pattern.compile(YOUTUBE_URL_PATTERN);
        Matcher matcher = pattern.matcher(url);

        // Ensure the URL matches and has the expected number of capture groups.
        if (!matcher.matches() || matcher.groupCount() != 1) {
            throw new YouTubeException("Video Id not found");
        }

        // Return the captured video ID wrapped in an Optional.
        return matcher.group(1);
    }

    /**
     * Get video info by URL
     *
     * @param url The YouTube video URL
     * @return VideoInfo object with video details
     * @throws YouTubeException if video ID is not found
     */
    VideoInfo getVideoInfo(String url);

    /**
     * Get video transcript by URL and language code
     *
     * @param url          The YouTube video URL
     * @param languageCode The language code to download transcript in
     * @return VideoTranscript object with video transcript
     * @throws YouTubeException if video ID is not found
     */
    VideoTranscript getTranscript(String url, String languageCode);
}
