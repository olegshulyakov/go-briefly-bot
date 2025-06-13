package com.github.youtubebriefly.service;

import com.github.youtubebriefly.dao.VideoSummaryRepository;
import com.github.youtubebriefly.exception.YouTubeException;
import com.github.youtubebriefly.model.VideoInfo;
import com.github.youtubebriefly.model.VideoSummary;
import com.github.youtubebriefly.model.VideoTranscript;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.ai.openai.api.common.OpenAiApiClientErrorException;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.time.LocalDateTime;

/**
 * Service responsible for generating summaries of YouTube videos.
 */
@Service
@RequiredArgsConstructor
public class YouTubeSummaryService {

    private static final Logger logger = LoggerFactory.getLogger(YouTubeSummaryService.class);

    private final VideoSummaryRepository videoSummaryRepository;
    private final YouTubeService youtubeService;
    private final SummaryService summaryService;

    /**
     * Retrieves a summary for a YouTube video.
     *
     * @param videoUrl The URL of the YouTube video.
     * @param language The desired language for the summary.
     * @return The VideoSummary object containing the summary information.
     * @throws YouTubeException              if the videoUrl is invalid
     * @throws YouTubeException              If there is an issue fetching video information from YouTube.
     * @throws OpenAiApiClientErrorException If there is an issue generating the summary using the OpenAI API.
     */
    public VideoSummary getSummary(String videoUrl, String language) throws YouTubeException, OpenAiApiClientErrorException {
        logger.info("Get video summary: {}-{}", language, videoUrl);

        if (!StringUtils.hasText(videoUrl)) {
            throw new IllegalArgumentException("Video URL cannot be null or empty.");
        }
        if (!StringUtils.hasText(language)) {
            throw new IllegalArgumentException("Language cannot be null or empty.");
        }

        String videoId = YouTubeService.getVideoId(videoUrl);

        if (videoSummaryRepository.existsByTypeAndVideoIdAndLanguage("youtube", videoId, language)) {
            logger.debug("Using video summary from cache: {}-{}-{}", "youtube", videoId, language);
            return videoSummaryRepository.findByTypeAndVideoIdAndLanguage("youtube", videoId, language);
        }

        // Step 1: Fetch video info
        VideoInfo videoInfo = youtubeService.getVideoInfo(videoUrl);

        // Step 2: Fetch transcript
        VideoTranscript videoTranscript = youtubeService.getTranscript(videoUrl, videoInfo.getLanguage());

        // Step 3: Summarize transcript
        String summary = summaryService.generateSummary(videoTranscript.getTranscript(), language);

        VideoSummary videoSummary = new VideoSummary("youtube", videoId, language, videoInfo.getTitle(), summary.trim(), LocalDateTime.now());
        videoSummaryRepository.save(videoSummary);

        return videoSummary;
    }
}
