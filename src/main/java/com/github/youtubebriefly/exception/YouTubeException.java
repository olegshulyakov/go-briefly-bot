package com.github.youtubebriefly.exception;

/**
 * A custom exception class specifically for YouTube-related errors.
 * <p>
 * This exception extends {@link RuntimeException}, indicating that it represents
 * unrecoverable errors that should be handled by the calling code.  It provides
 * constructors for passing a message and an optional cause (another exception
 * that led to this exception).
 * </p>
 */
public class YouTubeException extends RuntimeException {

    /**
     * Constructs a {@code YouTubeException} with the specified message.
     *
     * @param message The error message.  Should be descriptive and informative.
     */
    public YouTubeException(String message) {
        super(message);
    }

    /**
     * Constructs a {@code YouTubeException} with the specified message and cause.
     *
     * @param message The error message.
     * @param cause The underlying cause of this exception.  Can be {@code null}.
     */
    public YouTubeException(String message, Throwable cause) {
        super(message, cause);
    }
}
