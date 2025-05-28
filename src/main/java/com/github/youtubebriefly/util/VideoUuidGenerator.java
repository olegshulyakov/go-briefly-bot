package com.github.youtubebriefly.util;

public interface VideoUuidGenerator {
    /**
     * Generates a UUID based on the type and ID of the video.
     *
     * @param type Video type
     * @param id   Video ID
     * @return Generated UUID
     */
    static String getUuid(String type, String id) {
        return type + "-" + id;
    }
}
