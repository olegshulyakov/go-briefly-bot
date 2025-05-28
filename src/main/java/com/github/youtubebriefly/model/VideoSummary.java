package com.github.youtubebriefly.model;

import com.github.youtubebriefly.util.VideoUuidGenerator;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Entity
@Table(name = "VideoSummary")
@Data
@NoArgsConstructor
public class VideoSummary {
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
    private String title;
    private String summary;
    private String language;
    private LocalDateTime createdAt;

    public VideoSummary(String type, String videoId, String title, String summary, String language, LocalDateTime createdAt) {
        this.uuid = VideoUuidGenerator.getUuid(type, videoId);
        this.type = type;
        this.videoId = videoId;
        this.title = title;
        this.summary = summary;
        this.language = language;
        this.createdAt = createdAt;
    }
}
