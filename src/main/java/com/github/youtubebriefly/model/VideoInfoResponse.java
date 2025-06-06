package com.github.youtubebriefly.model;

import java.time.Instant;
import java.time.ZoneOffset;

public record VideoInfoResponse(
        String type,
        String videoId,
        String language,
        String uploader,
        String title,
        String thumbnail,
        Instant createdAt
) {
    public VideoInfoResponse(VideoInfo videoInfo) {
        this(videoInfo.getType(), videoInfo.getVideoId(), videoInfo.getLanguage(), videoInfo.getUploader(), videoInfo.getTitle(), videoInfo.getThumbnail(), videoInfo.getCreatedAt().toInstant(ZoneOffset.UTC));
    }
}
