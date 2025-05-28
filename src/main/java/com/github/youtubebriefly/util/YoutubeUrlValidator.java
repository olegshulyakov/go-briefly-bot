package com.github.youtubebriefly.util;

import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * <p>
 * Interface defining methods for validating YouTube URLs and extracting video IDs.
 * Provides a default implementation for common validation and parsing tasks.
 * </p>
 *
 * <p>
 * This class exemplifies the use of Regular Expressions for URL validation and extraction,
 * and utilizes the Optional type for safe handling of potential null or missing values.
 * It adheres to the Single Responsibility Principle, keeping the validation and ID extraction
 * logic encapsulated within this interface.
 * </p>
 */
public interface YoutubeUrlValidator {

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
    default boolean isValidYoutubeUrl(String url) {
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
     * @return An {@link Optional<String>} containing the video ID if the URL is a valid YouTube video URL;
     * otherwise, an empty {@link Optional}. This avoids NullPointerExceptions and provides a cleaner way
     * to represent the possibility of a missing ID.
     * @throws IllegalArgumentException if the input URL is null.
     */
    default Optional<String> getYoutubeId(String url) {
        // Null check
        if (url == null) {
            return Optional.empty();
        }

        Pattern pattern = Pattern.compile(YOUTUBE_URL_PATTERN);
        Matcher matcher = pattern.matcher(url);

        // Ensure the URL matches and has the expected number of capture groups.
        if (!matcher.matches() || matcher.groupCount() != 1) {
            return Optional.empty();
        }

        // Return the captured video ID wrapped in an Optional.
        return Optional.of(matcher.group(1));
    }
}
