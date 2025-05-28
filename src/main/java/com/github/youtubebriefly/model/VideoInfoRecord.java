package com.github.youtubebriefly.model;

public record VideoInfoRecord(
        String type,
        String id,
        String uploader,
        String title,
        String thumbnail
) { }
