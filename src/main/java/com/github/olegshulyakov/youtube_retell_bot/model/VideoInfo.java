package com.github.olegshulyakov.youtube_retell_bot.model;

import com.github.olegshulyakov.youtube_retell_bot.util.VideoUuidGenerator;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;

@Entity
@Table(name = "VideoInfo")
@Data
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
    public VideoInfo() {
    }

    public VideoInfo(String uuid, String type, String id, String uploader, String title, String thumbnail) {
        this.uuid = uuid;
        this.type = type;
        this.id = id;
        this.uploader = uploader;
        this.title = title;
        this.thumbnail = thumbnail;
    }

    public VideoInfo(VideoInfoRecord videoInfoRecord) {
        this(VideoUuidGenerator.getUuid(videoInfoRecord.type(), videoInfoRecord.id()), videoInfoRecord.type(), videoInfoRecord.id(), videoInfoRecord.uploader(), videoInfoRecord.title(), videoInfoRecord.thumbnail());
    }

    public VideoInfoRecord toVideoInfo() {
        return new VideoInfoRecord(this.type, this.id, this.uploader, this.title, this.thumbnail);
    }
}

