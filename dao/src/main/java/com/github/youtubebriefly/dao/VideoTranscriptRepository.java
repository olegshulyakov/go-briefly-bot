package com.github.youtubebriefly.dao;

import com.github.youtubebriefly.model.VideoTranscript;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

/**
 * Spring Data JPA repository interface for managing {@link VideoTranscript} entities.
 */
@Repository
public interface VideoTranscriptRepository extends JpaRepository<VideoTranscript, String> {

    /**
     * Checks if a video transcript exists with the specified type, video ID, and language.
     *
     * @param type The type of the transcript (e.g., "auto", "manual").
     * @param videoId The ID of the video associated with the transcript.
     * @param language The language of the transcript (e.g., "en", "es").
     * @return {@code true} if a transcript with the specified criteria exists, {@code false} otherwise.
     */
    boolean existsByTypeAndVideoIdAndLanguage(String type, String videoId, String language);

    /**
     * Retrieves a video transcript with the specified type, video ID, and language.
     *
     * @param type The type of the transcript (e.g., "auto", "manual").
     * @param videoId The ID of the video associated with the transcript.
     * @param language The language of the transcript (e.g., "en", "es").
     * @return The {@link VideoTranscript} object if found, or {@code null} if not found.
     */
    VideoTranscript findByTypeAndVideoIdAndLanguage(String type, String videoId, String language);
}