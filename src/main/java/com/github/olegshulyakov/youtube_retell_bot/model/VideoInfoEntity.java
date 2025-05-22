package com.github.olegshulyakov.youtube_retell_bot.model;

import com.github.olegshulyakov.youtube_retell_bot.util.VideoUuidGenerator;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;

@Entity
@Table(name = "VideoInfo")
@Data
public class VideoInfoEntity {
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
    private String id;

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

    // Required for JPA
    public VideoInfoEntity() {
    }

    public VideoInfoEntity(String uuid, String type, String id, String uploader, String title, String thumbnail) {
        this.uuid = uuid;
        this.type = type;
        this.id = id;
        this.uploader = uploader;
        this.title = title;
        this.thumbnail = thumbnail;
    }

    public VideoInfoEntity(VideoInfo videoInfo) {
        this(VideoUuidGenerator.getUuid(videoInfo.type(), videoInfo.id()), videoInfo.type(), videoInfo.id(), videoInfo.uploader(), videoInfo.title(), videoInfo.thumbnail());
    }

    public VideoInfo toVideoInfo() {
        return new VideoInfo(this.type, this.id, this.uploader, this.title, this.thumbnail);
    }
}

