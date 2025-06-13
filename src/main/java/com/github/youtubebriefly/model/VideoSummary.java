package com.github.youtubebriefly.model;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

/**
 * Represents a summarized video, identified by its type, video ID, and language.
 * This entity stores the title and summary of a video, along with the creation timestamp.
 * <p>
 * The primary key is a composite key consisting of the video type, video ID, and language.
 */
@Entity
@IdClass(VideoKey.class)
@Table(name = "VideoSummary")
@Data
@NoArgsConstructor
public class VideoSummary {

    /**
     * The type of video platform (e.g., "youtube").
     */
    @Id
    private String type;

    /**
     * The unique video ID from the originating platform (e.g., YouTube video ID).
     */
    @Id
    private String videoId;

    /**
     * The language code of the summary (e.g., "en", "es", "fr").
     */
    @Id
    private String language;

    /**
     * The title of the video.
     */
    @Column(nullable = false)
    private String title;

    /**
     * The summary of the video content.
     */
    @Column(nullable = false)
    private String summary;

    /**
     * The timestamp indicating when the summary was created.
     */
    @Column(nullable = false)
    private LocalDateTime createdAt;

    /**
     * Constructs a VideoSummary object.
     *
     * @param type      The video type (e.g., "youtube").
     * @param videoId   The unique video ID.
     * @param language  The language of the summary.
     * @param title     The title of the video.
     * @param summary   The summary of the video.
     * @param createdAt The creation timestamp.
     */
    public VideoSummary(String type, String videoId, String language, String title, String summary, LocalDateTime createdAt) {
        this.type = type;
        this.videoId = videoId;
        this.language = language;
        this.title = title;
        this.summary = summary;
        this.createdAt = createdAt;
    }
}
