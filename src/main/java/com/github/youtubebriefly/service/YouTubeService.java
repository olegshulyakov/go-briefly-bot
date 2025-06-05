package com.github.youtubebriefly.service;

import com.github.youtubebriefly.exception.YouTubeException;
import com.github.youtubebriefly.model.VideoInfo;
import com.github.youtubebriefly.model.VideoTranscript;

/**
 * YouTube Service
 * Handles video info retrieval and transcript downloading.
 */
public interface YouTubeService {
    /**
     * Get video info by URL
     *
     * @param url The YouTube video URL
     * @return VideoInfo object with video details
     * @throws YouTubeException if video ID is not found
     */
    VideoInfo getVideoInfo(String url);

    /**
     * Get video transcript by URL and language code
     *
     * @param url          The YouTube video URL
     * @param languageCode The language code to download transcript in
     * @return VideoTranscript object with video transcript
     * @throws YouTubeException if video ID is not found
     */
    VideoTranscript getTranscript(String url, String languageCode);
}
