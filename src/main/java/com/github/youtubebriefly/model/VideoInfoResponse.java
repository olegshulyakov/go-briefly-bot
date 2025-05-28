package com.github.youtubebriefly.model;

public record VideoInfoResponse(
        String type,
        String videoId,
        String uploader,
        String title,
        String thumbnail
) {
    public VideoInfoResponse(VideoInfo videoInfo) {
        this(videoInfo.getType(), videoInfo.getVideoId(), videoInfo.getUploader(), videoInfo.getTitle(), videoInfo.getThumbnail());
    }
}
