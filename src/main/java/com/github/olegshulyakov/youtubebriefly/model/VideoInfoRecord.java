package com.github.olegshulyakov.youtubebriefly.model;

public record VideoInfoRecord(
        String type,
        String id,
        String uploader,
        String title,
        String thumbnail
) { }
