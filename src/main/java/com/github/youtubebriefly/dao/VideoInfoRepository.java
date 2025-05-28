package com.github.youtubebriefly.dao;

import com.github.youtubebriefly.model.VideoInfo;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface VideoInfoRepository extends JpaRepository<VideoInfo, String> {
    boolean existsByTypeAndVideoId(String type, String videoId);

    VideoInfo findByTypeAndVideoId(String type, String videoId);
}