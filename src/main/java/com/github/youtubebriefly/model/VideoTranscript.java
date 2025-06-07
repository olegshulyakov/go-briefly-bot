package com.github.youtubebriefly.model;

import jakarta.persistence.*;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

/**
 * Represents a video transcript stored in the system.  This entity
 * handles transcripts for various video types (e.g., YouTube, Vimeo).
 * It uses a composite key consisting of video type, video ID, and language
 * to uniquely identify each transcript.
 */
@Entity
@IdClass(VideoKey.class)
@Table(name = "VideoTranscript")
@Data
@NoArgsConstructor
public class VideoTranscript {

    /**
     * Video type (e.g., "youtube").  This identifies the source of the video.
     */
    @Id
    private String type;

    /**
     * Video unique ID from the video platform API (e.g., YouTube API).
     */
    @Id
    private String videoId;

    /**
     * Video transcript language code (e.g., "en", "es", "fr").
     */
    @Id
    private String language;

    /**
     * The actual video transcript text.
     */
    @Column(nullable = false)
    private String transcript;

    /**
     * Timestamp indicating when the transcript was created.  This field is mandatory.
     */
    @Column(nullable = false)
    private LocalDateTime createdAt;


    /**
     * Constructs a new {@code VideoTranscript} object.
     *
     * @param type       The video type (e.g., "youtube").  Must not be null or empty.
     * @param videoId    The unique video ID. Must not be null or empty.
     * @param language   The language code for the transcript.  Must not be null or empty.
     * @param createdAt  The timestamp when the transcript was created. Must not be null.
     * @param transcript The transcript text. Must not be null or empty.
     * @throws IllegalArgumentException if any of the required parameters are null or empty.
     */
    public VideoTranscript(String type, String videoId, String language, LocalDateTime createdAt, String transcript) {
        this.type = type;
        this.videoId = videoId;
        this.language = language;
        this.createdAt = createdAt;
        this.transcript = transcript;
    }
}
