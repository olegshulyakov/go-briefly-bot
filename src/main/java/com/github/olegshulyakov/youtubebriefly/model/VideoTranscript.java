package com.github.olegshulyakov.youtubebriefly.model;

import com.github.olegshulyakov.youtubebriefly.util.VideoUuidGenerator;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import lombok.Data;

@Entity
@Table(name = "VideoTranscript")
@Data
public class VideoTranscript {
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
    public VideoTranscript() {
    }

    public VideoTranscript(String uuid, String type, String id, String transcript) {
        this.uuid = uuid;
        this.type = type;
        this.id = id;
        this.transcript = transcript;
    }

    public VideoTranscript(VideoTranscriptRecord videoTranscriptRecord) {
        this(VideoUuidGenerator.getUuid(videoTranscriptRecord.type(), videoTranscriptRecord.id()), videoTranscriptRecord.type(), videoTranscriptRecord.id(), videoTranscriptRecord.transcript());
    }

    public VideoTranscriptRecord toVideoTranscript() {
        return new VideoTranscriptRecord(this.type, this.id, this.transcript);
    }
}
