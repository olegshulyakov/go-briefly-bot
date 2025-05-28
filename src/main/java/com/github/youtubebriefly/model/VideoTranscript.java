package com.github.youtubebriefly.model;

import com.github.youtubebriefly.util.VideoUuidGenerator;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;
import lombok.NoArgsConstructor;

@Entity
@Table(name = "VideoTranscript")
@Data
@NoArgsConstructor
public class VideoTranscript {
    /**
     * Unique identifier for the video.
     */
    @Id
    private String uuid;
    /**
     * Video type (e.g., "youtube").
     */
    private String type;
    /**
     * Video unique ID from the YouTube API.
     */
    private String videoId;
    /**
     * Video title.
     */
    private String transcript;

    public VideoTranscript(String type, String videoId, String transcript) {
        this.uuid = VideoUuidGenerator.getUuid(type, videoId);
        this.type = type;
        this.videoId = videoId;
        this.transcript = transcript;
    }
}
