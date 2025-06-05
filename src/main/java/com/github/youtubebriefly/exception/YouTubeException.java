package com.github.youtubebriefly.exception;

public class YouTubeException extends RuntimeException {
    public YouTubeException(String message) {
        super(message);
    }

    public YouTubeException(String message, Throwable cause) {
        super(message, cause);
    }
}
