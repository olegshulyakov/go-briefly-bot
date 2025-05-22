package com.github.olegshulyakov.youtube_retell_bot.model;

public record VideoTranscript(
        String type,
        String id,
        String transcript
) { }
