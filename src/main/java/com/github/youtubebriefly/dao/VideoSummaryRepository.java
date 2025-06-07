package com.github.youtubebriefly.dao;

import com.github.youtubebriefly.model.VideoSummary;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface VideoSummaryRepository extends JpaRepository<VideoSummary, String> {
    boolean existsByTypeAndVideoIdAndLanguage(String type, String videoId, String language);

    VideoSummary findByTypeAndVideoIdAndLanguage(String type, String videoId, String language);
}
