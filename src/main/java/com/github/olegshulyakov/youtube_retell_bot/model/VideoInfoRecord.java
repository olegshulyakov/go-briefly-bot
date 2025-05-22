package com.github.olegshulyakov.youtube_retell_bot.model;

public record VideoInfoRecord(
        String type,
        String id,
        String uploader,
        String title,
        String thumbnail
) { }
