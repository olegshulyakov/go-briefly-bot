package com.github.youtubebriefly.model;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Entity
@IdClass(VideoKey.class)
@Table(name = "VideoTranscript")
@Data
@NoArgsConstructor
public class VideoTranscript {
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
     * Video transcript language code.
     */
    @Id
    private String language;
    /**
     * Video transcript.
     */
    @Column(nullable = false)
    private String transcript;

    @Column(nullable = false)
    private LocalDateTime createdAt;

    public VideoTranscript(String type, String videoId, String language, LocalDateTime createdAt, String transcript) {
        this.type = type;
        this.videoId = videoId;
        this.language = language;
        this.createdAt = createdAt;
        this.transcript = transcript;
    }
}
