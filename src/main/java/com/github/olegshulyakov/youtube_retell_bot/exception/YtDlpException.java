package com.github.olegshulyakov.youtube_retell_bot.exception;

public class YtDlpException extends RuntimeException {
    public YtDlpException(String message) {
        super(message);
    }

    public YtDlpException(String message, Throwable cause) {
        super(message, cause);
    }
}
