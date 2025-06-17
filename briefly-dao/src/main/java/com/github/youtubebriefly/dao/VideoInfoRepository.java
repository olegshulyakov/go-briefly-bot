package com.github.youtubebriefly.dao;

import com.github.youtubebriefly.model.VideoInfo;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

/**
 * Repository interface for managing {@link VideoInfo} entities.
 *
 * @author [Your Name/Team Name]
 */
@Repository
public interface VideoInfoRepository extends JpaRepository<VideoInfo, String> {

    /**
     * Checks if a video information entry exists for the given type and video ID.
     *
     * @param type    The type of the video information.  Must not be null or empty.
     * @param videoId The ID of the video. Must not be null or empty.
     * @return {@code true} if an entry exists with the specified type and video ID,
     *         {@code false} otherwise.  Returns false if either `type` or `videoId` is null or empty.
     * @throws IllegalArgumentException if type or videoId is null or empty
     */
    boolean existsByTypeAndVideoId(String type, String videoId);

    /**
     * Retrieves a {@link VideoInfo} entity based on the specified type and video ID.
     *
     * @param type    The type of the video information.  Must not be null or empty.
     * @param videoId The ID of the video. Must not be null or empty.
     * @return The {@link VideoInfo} entity if found, or {@code null} if no entry exists
     *         with the specified type and video ID.  Returns null if either `type` or `videoId` is null or empty.
     * @throws IllegalArgumentException if type or videoId is null or empty
     */
    VideoInfo findByTypeAndVideoId(String type, String videoId);
}