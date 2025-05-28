package com.github.youtubebriefly.dao;

import com.github.youtubebriefly.model.VideoTranscript;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface VideoTranscriptRepository extends JpaRepository<VideoTranscript, String> {
    boolean existsByTypeAndVideoIdAndLanguageCode(String type, String videoId, String languageCode);

    VideoTranscript findByTypeAndVideoIdAndLanguageCode(String type, String videoId, String languageCode);
}