package com.github.youtubebriefly.controller;

import com.github.youtubebriefly.model.VideoInfo;
import com.github.youtubebriefly.model.VideoInfoResponse;
import com.github.youtubebriefly.model.VideoTranscript;
import com.github.youtubebriefly.model.VideoTranscriptResponse;
import com.github.youtubebriefly.service.YouTubeService;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

/**
 * REST controller for handling video information and transcript requests.
 *
 * This controller exposes endpoints to retrieve video metadata and transcripts
 * from YouTube videos using the {@link YouTubeService}.
 */
@RestController
@RequestMapping("/video")
@RequiredArgsConstructor
public class VideoController {

    private static final Logger logger = LoggerFactory.getLogger(VideoController.class);

    /**
     * The YouTube service used to fetch video information and transcripts.
     */
    private final YouTubeService youTubeService;

    /**
     * Retrieves information about a YouTube video.
     *
     * @param url The URL of the YouTube video.
     * @return A {@link ResponseEntity} containing the {@link VideoInfoResponse} with video information,
     *         or an appropriate error response if the video information cannot be retrieved.
     */
    @GetMapping("/info")
    public ResponseEntity<VideoInfoResponse> getVideoInfo(@RequestParam String url) {
        logger.info("Received request to get video info for URL: {}", url);
        VideoInfo videoInfo = youTubeService.getVideoInfo(url);
        return ResponseEntity.ok(new VideoInfoResponse(videoInfo));
    }

    /**
     * Retrieves the transcript of a YouTube video.
     *
     * @param url          The URL of the YouTube video.
     * @param languageCode The language code for the transcript (e.g., "en", "es").  Defaults to "en".
     * @return A {@link ResponseEntity} containing the {@link VideoTranscriptResponse} with the video transcript,
     *         or an appropriate error response if the transcript cannot be retrieved.
     */
    @GetMapping("/transcript")
    public ResponseEntity<VideoTranscriptResponse> getTranscript(@RequestParam String url, @RequestParam(required = false, defaultValue = "en") String languageCode) {
        logger.info("Received request to get transcript for URL: {} and language code: {}", url, languageCode);
        VideoTranscript transcript = youTubeService.getTranscript(url, languageCode);
        return ResponseEntity.ok(new VideoTranscriptResponse(transcript));
    }
}
