package com.github.olegshulyakov.youtube_retell_bot.dao;

import com.github.olegshulyakov.youtube_retell_bot.model.VideoInfoEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface VideoInfoEntityRepository extends JpaRepository<VideoInfoEntity, String> {
    boolean existsByUuid(String uuid);

    VideoInfoEntity findByUuid(String uuid);
}