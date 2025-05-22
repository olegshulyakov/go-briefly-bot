package com.github.olegshulyakov.youtube_retell_bot.dao;

import com.github.olegshulyakov.youtube_retell_bot.model.VideoTranscriptEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface VideoTranscriptEntityRepository extends JpaRepository<VideoTranscriptEntity, String> {
    boolean existsByUuid(String uuid);

    VideoTranscriptEntity findByUuid(String uuid);
}