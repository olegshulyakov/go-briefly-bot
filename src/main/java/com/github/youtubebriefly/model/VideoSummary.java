package com.github.youtubebriefly.model;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Entity
@IdClass(VideoKey.class)
@Table(name = "VideoSummary")
@Data
@NoArgsConstructor
public class VideoSummary {
    /**
     * Video type (e.g., "youtube").
     */
    @Id
    private String type;
    /**
     * Video unique ID from the YouTube API.
     */
    @Id
    private String videoId;
    /**
     * Summary language code.
     */
    @Id
    private String languageCode;

    @Column(nullable = false)
    private String title;

    @Column(nullable = false)
    private String summary;

    @Column(nullable = false)
    private LocalDateTime createdAt;

    public VideoSummary(String type, String videoId, String title, String summary, String languageCode, LocalDateTime createdAt) {
        this.type = type;
        this.videoId = videoId;
        this.languageCode = languageCode;
        this.title = title;
        this.summary = summary;
        this.createdAt = createdAt;
    }
}
