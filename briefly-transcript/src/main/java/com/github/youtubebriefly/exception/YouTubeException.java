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
     * {@inheritDoc}
     */
    public YouTubeException(String message) {
        super(message);
    }

    /**
     * {@inheritDoc}
     */
    public YouTubeException(Throwable cause) {
        super(cause);
    }

    /**
     * {@inheritDoc}
     */
    public YouTubeException(String message, Throwable cause) {
        super(message, cause);
    }
}
