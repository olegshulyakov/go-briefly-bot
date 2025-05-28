package com.github.youtubebriefly.model;

import com.github.youtubebriefly.util.VideoUuidGenerator;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;
import lombok.NoArgsConstructor;

@Entity
@Table(name = "VideoInfo")
@Data
@NoArgsConstructor
public class VideoInfo {
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
     * Uploader's username.
     */
    private String uploader;

    /**
     * Video title.
     */
    private String title;

    /**
     * Video thumbnail URL.
     */
    private String thumbnail;

    public VideoInfo(String type, String videoId, String uploader, String title, String thumbnail) {
        this.uuid = VideoUuidGenerator.getUuid(type, videoId);
        this.type = type;
        this.videoId = videoId;
        this.uploader = uploader;
        this.title = title;
        this.thumbnail = thumbnail;
    }
}
