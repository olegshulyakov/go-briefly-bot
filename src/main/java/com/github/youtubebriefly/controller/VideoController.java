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

@RestController
@RequestMapping("/video")
@RequiredArgsConstructor
public class VideoController {
    private static final Logger logger = LoggerFactory.getLogger(VideoController.class);

    private final YouTubeService youTubeService;

    @GetMapping("/info")
    public ResponseEntity<VideoInfoResponse> getVideoInfo(@RequestParam String url) {
        VideoInfo videoInfo = youTubeService.getVideoInfo(url);
        return ResponseEntity.ok(new VideoInfoResponse(videoInfo));
    }

    @GetMapping("/transcript")
    public ResponseEntity<VideoTranscriptResponse> getTranscript(@RequestParam String url, @RequestParam(required = false, defaultValue = "en") String languageCode) {
        VideoTranscript transcript = youTubeService.getTranscript(url, languageCode);
        return ResponseEntity.ok(new VideoTranscriptResponse(transcript));
    }
}
