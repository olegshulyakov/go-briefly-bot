package com.github.youtubebriefly.model;

public record VideoTranscriptResponse(
        String type,
        String videoId,
        String transcript
) {
    public VideoTranscriptResponse(VideoTranscript videoTranscript) {
        this(videoTranscript.getType(), videoTranscript.getVideoId(), videoTranscript.getTranscript());
    }
}
