package com.github.olegshulyakov.youtubebriefly.dao;

import com.github.olegshulyakov.youtubebriefly.model.VideoTranscript;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface VideoTranscriptRepository extends JpaRepository<VideoTranscript, String> {
    boolean existsByUuid(String uuid);

    VideoTranscript findByUuid(String uuid);
}