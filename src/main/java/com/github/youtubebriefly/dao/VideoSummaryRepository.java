package com.github.youtubebriefly.dao;

import com.github.youtubebriefly.model.VideoSummary;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

/**
 * Repository interface for managing {@link VideoSummary} entities.
 */
@Repository
public interface VideoSummaryRepository extends JpaRepository<VideoSummary, String> {

    /**
     * Checks if a video summary exists with the specified type, video ID, and language.
     *
     * @param type The type of the video summary (e.g., "transcript", "keywords"). Must not be null or empty.
     * @param videoId The ID of the video. Must not be null or empty.
     * @param language The language of the video summary (e.g., "en", "es"). Must not be null or empty.
     * @return {@code true} if a video summary exists with the specified criteria, {@code false} otherwise.
     * @throws IllegalArgumentException if any of the input parameters are null or empty.
     */
    boolean existsByTypeAndVideoIdAndLanguage(String type, String videoId, String language);

    /**
     * Retrieves a video summary with the specified type, video ID, and language.
     *
     * @param type The type of the video summary (e.g., "transcript", "keywords"). Must not be null or empty.
     * @param videoId The ID of the video. Must not be null or empty.
     * @param language The language of the video summary (e.g., "en", "es"). Must not be null or empty.
     * @return The {@link VideoSummary} object if found, or {@code null} if no matching summary exists.
     * @throws IllegalArgumentException if any of the input parameters are null or empty.
     */
    VideoSummary findByTypeAndVideoIdAndLanguage(String type, String videoId, String language);
}
