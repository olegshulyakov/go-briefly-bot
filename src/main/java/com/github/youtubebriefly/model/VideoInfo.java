package com.github.youtubebriefly.model;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Entity
@IdClass(VideoKey.class)
@Table(name = "VideoInfo")
@Data
@NoArgsConstructor
public class VideoInfo {

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
     * Language code.
     */
    @Id
    private String language = "";

    /**
     * Uploader's username.
     */
    @Column(nullable = false)
    private String uploader;

    /**
     * Video title.
     */
    @Column(nullable = false)
    private String title;

    /**
     * Video thumbnail URL.
     */
    @Column(nullable = false)
    private String thumbnail;

    @Column(nullable = false)
    private LocalDateTime createdAt;

    public VideoInfo(String type, String videoId, String language, String uploader, String title, String thumbnail, LocalDateTime createdAt) {
        this.type = type;
        this.language = language;
        this.videoId = videoId;
        this.uploader = uploader;
        this.title = title;
        this.thumbnail = thumbnail;
        this.createdAt = createdAt;
    }
}
