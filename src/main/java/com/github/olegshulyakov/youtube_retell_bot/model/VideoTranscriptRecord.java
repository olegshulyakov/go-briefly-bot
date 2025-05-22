package com.github.olegshulyakov.youtube_retell_bot.model;

public record VideoTranscriptRecord(
        String type,
        String id,
        String transcript
) { }
