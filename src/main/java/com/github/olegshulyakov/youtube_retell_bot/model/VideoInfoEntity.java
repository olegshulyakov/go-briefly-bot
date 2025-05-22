package com.github.olegshulyakov.youtube_retell_bot.model;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;

@Entity
@Table(name = "VideoInfo")
@Data
public class VideoInfoEntity {
    @Id
    private String uuid;
    private String type;
    private String id;
    private String uploader;
    private String title;
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
        this(getUuid(videoInfo.type(), videoInfo.id()), videoInfo.type(), videoInfo.id(), videoInfo.uploader(), videoInfo.title(), videoInfo.thumbnail());
    }

    public VideoInfo toVideoInfo() {
        return new VideoInfo(this.type, this.id, this.uploader, this.title, this.thumbnail);
    }

    public static String getUuid(String type, String id) {
        return type + "-" + id;
    }
}
