package com.github.youtubebriefly.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.youtubebriefly.config.YouTubeConfig;
import com.github.youtubebriefly.dao.VideoInfoRepository;
import com.github.youtubebriefly.dao.VideoTranscriptRepository;
import com.github.youtubebriefly.exception.YouTubeException;
import com.github.youtubebriefly.file.SubtitleFormats;
import com.github.youtubebriefly.file.TranscriptFiles;
import com.github.youtubebriefly.file.Transcripts;
import com.github.youtubebriefly.model.VideoInfo;
import com.github.youtubebriefly.model.VideoTranscript;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.retry.annotation.Retryable;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;

/**
 * Service handles video info retrieval and transcript downloading using yt-dlp
 */
@Service
@RequiredArgsConstructor
public class YtDlpService implements YouTubeService {
    private static final Logger logger = LoggerFactory.getLogger(YtDlpService.class);
    private static final String SUBTITLE_FORMAT = SubtitleFormats.SRT.toString();

    private final VideoInfoRepository videoInfoRepository;
    private final VideoTranscriptRepository videoTranscriptRepository;
    private final YouTubeConfig youtubeConfig;

    /**
     * {@inheritDoc}
     */
    @Override
    @Retryable(retryFor = YouTubeException.class, maxAttempts = 3)
    public VideoInfo getVideoInfo(String url) {
        logger.info("Get video info: {}", url);
        String videoId = YouTubeService.getVideoId(url);

        if (videoInfoRepository.existsByTypeAndVideoId("youtube", videoId)) {
            logger.debug("Using video info from cache: {}-{}", "youtube", videoId);
            return videoInfoRepository.findByTypeAndVideoId("youtube", videoId);
        }

        List<String> args = Arrays.asList("--dump-json", url);
        String output = execYtDlpCommand(args);
        logger.debug("Got video info: {}-{}", "youtube", videoId);

        try {
            @SuppressWarnings("unchecked")
            Map<String, Object> jsonMap = new ObjectMapper().readValue(output, Map.class);
            VideoInfo videoInfo = new VideoInfo(
                    "youtube",
                    videoId,
                    (String) jsonMap.get("language"),
                    (String) jsonMap.get("uploader"),
                    (String) jsonMap.get("title"),
                    (String) jsonMap.get("thumbnail"),
                    LocalDateTime.now()
            );

            videoInfoRepository.save(videoInfo);

            return videoInfo;
        } catch (JsonProcessingException e) {
            throw new YouTubeException("Failed to parse yt-dlp output", e);
        }
    }

    /**
     * {@inheritDoc}
     */
    @Override
    @Retryable(retryFor = YouTubeException.class, maxAttempts = 3)
    public VideoTranscript getTranscript(String url, String language) {
        logger.info("Get video transcript: {}-{}", language, url);
        String videoId = YouTubeService.getVideoId(url);

        if (videoTranscriptRepository.existsByTypeAndVideoIdAndLanguage("youtube", videoId, language)) {
            logger.debug("Using video transcript from cache: {}-{}-{}", "youtube", videoId, language);
            return videoTranscriptRepository.findByTypeAndVideoIdAndLanguage("youtube", videoId, language);
        }

        List<String> args = Arrays.asList(
                "--no-progress",
                "--skip-download",
                "--write-subs",
                "--write-auto-subs",
                "--convert-subs",
                SUBTITLE_FORMAT,
                "--sub-lang",
                String.format("%s,%s_auto,-live_chat", language, language),
                "--output",
                String.format("subtitles_%s", videoId),
                url
        );
        execYtDlpCommand(args);
        logger.debug("Got video transcript: {}-{}-{}", "youtube", videoId, language);

        // Read the transcript file
        String transcript;
        try {
            String fileName = String.format("subtitles_%s.%s.%s", videoId, language, SUBTITLE_FORMAT);
            transcript = TranscriptFiles.readAndDelete(fileName);
        } catch (IOException e) {
            logger.error("Failed to read transcript file", e);
            throw new YouTubeException(e);
        }

        VideoTranscript videoTranscript = new VideoTranscript("youtube", videoId, language, LocalDateTime.now(), Transcripts.cleanSRT(transcript));
        videoTranscriptRepository.save(videoTranscript);

        return videoTranscript;
    }

    /**
     * Execute yt-dlp command to get video info
     *
     * @param args Command arguments
     * @return Output from yt-dlp command
     * @throws YouTubeException if command fails
     */
    private String execYtDlpCommand(List<String> args) {
        List<String> command = new ArrayList<>(args.size() + 3);
        command.add("yt-dlp");

        if (StringUtils.hasText(youtubeConfig.getYtDlpProxy())) {
            command.add("--proxy");
            command.add(youtubeConfig.getYtDlpProxy());
        }
        command.addAll(args);

        logger.debug("Executing command: {}", String.join(" ", command));

        Process process;
        StringBuilder output = new StringBuilder();
        StringBuilder errorOutput = new StringBuilder();

        try {
            process = new ProcessBuilder(command).start();

            try (BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()));
                 BufferedReader errorReader = new BufferedReader(new InputStreamReader(process.getErrorStream()))) {

                String line;
                while ((line = reader.readLine()) != null) {
                    output.append(line).append("\n");
                }

                String errorLine;
                while ((errorLine = errorReader.readLine()) != null) {
                    errorOutput.append(errorLine).append("\n");
                }
            }

            process.waitFor(30, TimeUnit.SECONDS);
        } catch (IOException | InterruptedException e) {
            logger.error("Failed to read video info", e);
            throw new YouTubeException("Failed to read video info", e);
        }

        int exitCode = process.exitValue();
        if (exitCode != 0) {
            logger.warn("yt-dlp finished with exit code {}\n{}", exitCode, errorOutput);
            throw new YouTubeException("yt-dlp command failed with exit code " + exitCode);
        }
        return output.toString();
    }
}
