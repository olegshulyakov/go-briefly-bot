package com.github.olegshulyakov.youtube_retell_bot.dao;

import com.github.olegshulyakov.youtube_retell_bot.model.VideoTranscript;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface VideoTranscriptRepository extends JpaRepository<VideoTranscript, String> {
    boolean existsByUuid(String uuid);

    VideoTranscript findByUuid(String uuid);
}