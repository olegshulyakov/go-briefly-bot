package com.github.olegshulyakov.youtube_retell_bot.dao;

import com.github.olegshulyakov.youtube_retell_bot.model.VideoInfo;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface VideoInfoRepository extends JpaRepository<VideoInfo, String> {
    boolean existsByUuid(String uuid);

    VideoInfo findByUuid(String uuid);
}