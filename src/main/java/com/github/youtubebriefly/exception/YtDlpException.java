package com.github.youtubebriefly.exception;

public class YtDlpException extends RuntimeException {
    public YtDlpException(String message) {
        super(message);
    }

    public YtDlpException(String message, Throwable cause) {
        super(message, cause);
    }
}
