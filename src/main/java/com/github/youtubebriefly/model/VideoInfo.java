package com.github.youtubebriefly.model;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

/**
 * Represents video information, likely sourced from an external platform like YouTube.
 * This entity is designed to store metadata about videos, uniquely identified by video type,
 * video ID, and language.  The composite key ensures uniqueness across these dimensions.
 */
@Entity
@IdClass(VideoKey.class)
@Table(name = "VideoInfo")
@Data
@NoArgsConstructor
public class VideoInfo {

    /**
     * Video type (e.g., "youtube").  This is part of the composite primary key.
     */
    @Id
    private String type;

    /**
     * Video unique ID from the source platform (e.g., YouTube API).
     * This is part of the composite primary key.
     */
    @Id
    private String videoId;

    /**
     * Language code (e.g., "en", "es", "fr").  This is part of the composite primary key.
     * Defaults to an empty string if not provided.
     */
    @Id
    private String language = "";

    /**
     * Uploader's username or channel name.  Cannot be null.
     */
    @Column(nullable = false)
    private String uploader;

    /**
     * Video title.  Cannot be null.
     */
    @Column(nullable = false)
    private String title;

    /**
     * Video thumbnail URL.  Cannot be null.
     */
    @Column(nullable = false)
    private String thumbnail;

    /**
     * Timestamp indicating when the video information was first created/recorded.
     * Cannot be null.
     */
    @Column(nullable = false)
    private LocalDateTime createdAt;

    /**
     * Constructs a {@link VideoInfo} object with all required fields.
     *
     * @param type      The video type (e.g., "youtube").
     * @param videoId   The unique video ID from the source platform.
     * @param language  The language code of the video.
     * @param uploader  The uploader's username or channel name.
     * @param title     The video title.
     * @param thumbnail The URL of the video thumbnail.
     * @param createdAt The timestamp when the video information was created.
     */
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
