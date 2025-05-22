package com.github.olegshulyakov.youtube_retell_bot.model;

import com.github.olegshulyakov.youtube_retell_bot.util.VideoUuidGenerator;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;

@Entity
@Table(name = "VideoTranscript")
@Data
public class VideoTranscriptEntity {
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
     * Video title.
     */
    private String transcript;

    // Required for JPA
    public VideoTranscriptEntity() {
    }

    public VideoTranscriptEntity(String uuid, String type, String id, String transcript) {
        this.uuid = uuid;
        this.type = type;
        this.id = id;
        this.transcript = transcript;
    }

    public VideoTranscriptEntity(VideoTranscript videoTranscript) {
        this(VideoUuidGenerator.getUuid(videoTranscript.type(), videoTranscript.id()), videoTranscript.type(), videoTranscript.id(), videoTranscript.transcript());
    }

    public VideoTranscript toVideoTranscript() {
        return new VideoTranscript(this.type, this.id, this.transcript);
    }
}
