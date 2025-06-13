package com.github.youtubebriefly.model;

import java.time.Instant;
import java.time.ZoneOffset;

/**
 * Represents the response containing video information. This record is designed to be immutable
 * and encapsulates key details about a video. It is commonly used as a DTO (Data Transfer Object)
 * in microservice communication.
 */
public record VideoInfoResponse(
        String type,
        String videoId,
        String language,
        String uploader,
        String title,
        String thumbnail,
        Instant createdAt
) {
    /**
     * Constructs a {@code VideoInfoResponse} from a {@code VideoInfo} object.
     *
     * @param videoInfo The source {@code VideoInfo} object.
     */
    public VideoInfoResponse(VideoInfo videoInfo) {
        this(videoInfo.getType(), videoInfo.getVideoId(), videoInfo.getLanguage(), videoInfo.getUploader(), videoInfo.getTitle(), videoInfo.getThumbnail(), videoInfo.getCreatedAt().toInstant(ZoneOffset.UTC));
    }
}
