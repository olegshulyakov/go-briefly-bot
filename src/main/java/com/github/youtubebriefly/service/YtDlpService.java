package com.github.youtubebriefly.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.youtubebriefly.dao.VideoInfoRepository;
import com.github.youtubebriefly.dao.VideoTranscriptRepository;
import com.github.youtubebriefly.exception.YouTubeException;
import com.github.youtubebriefly.model.VideoInfo;
import com.github.youtubebriefly.model.VideoTranscript;
import com.github.youtubebriefly.util.TranscriptCleaner;
import com.github.youtubebriefly.util.YoutubeUrlValidator;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.io.*;
import java.time.LocalDateTime;
import java.util.*;
import java.util.concurrent.TimeUnit;

/**
 * Controller for YouTube Service
 * Handles video info retrieval and transcript downloading using yt-dlp
 */
@Service
@RequiredArgsConstructor
public class YtDlpService implements YouTubeService, YoutubeUrlValidator, TranscriptCleaner {
    private static final Logger logger = LoggerFactory.getLogger(YtDlpService.class);

    private final VideoInfoRepository videoInfoRepository;
    private final VideoTranscriptRepository videoTranscriptRepository;
    private final String ytDlpProxy;

    /**
     * {@inheritDoc}
     */
    @Override
    public VideoInfo getVideoInfo(String url) {
        logger.info("Get video info: {}", url);
        Optional<String> oVideoId = getYoutubeId(url);
        if (oVideoId.isEmpty()) {
            throw new YouTubeException("Video Id not found");
        }
        String videoId = oVideoId.get();

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
    public VideoTranscript getTranscript(String url, String languageCode) {
        logger.info("Get video transcript: {}", url);
        Optional<String> oVideoId = getYoutubeId(url);
        if (oVideoId.isEmpty()) {
            throw new YouTubeException("Video Id not found");
        }
        String videoId = oVideoId.get();

        if (videoTranscriptRepository.existsByTypeAndVideoIdAndLanguageCode("youtube", videoId, languageCode)) {
            logger.debug("Using video transcript from cache: {}-{}-{}", "youtube", videoId, languageCode);
            return videoTranscriptRepository.findByTypeAndVideoIdAndLanguageCode("youtube", videoId, languageCode);
        }

        List<String> args = Arrays.asList(
                "--no-progress",
                "--skip-download",
                "--write-subs",
                "--write-auto-subs",
                "--convert-subs",
                "srt",
                "--sub-lang",
                String.format("%s,%s_auto,-live_chat", languageCode, languageCode),
                "--output",
                String.format("subtitles_%s", videoId),
                url
        );
        execYtDlpCommand(args);
        logger.debug("Got video transcript: {}-{}-{}", "youtube", videoId, languageCode);

        // Generate the file name
        String fileName = String.format("subtitles_%s.%s.srt", videoId, languageCode);
        File transcriptFile = new File(fileName);

        // Check if file has been created
        if (!transcriptFile.exists() || !transcriptFile.isFile()) {
            logger.error("Transcript file not found: {}", transcriptFile.getAbsolutePath());
            throw new YouTubeException("Transcript file not found");
        }

        // Read the transcript file
        StringBuilder transcript = new StringBuilder();
        try (BufferedReader reader = new BufferedReader(new FileReader(transcriptFile))) {
            String line;
            while ((line = reader.readLine()) != null) {
                transcript.append(line).append("\n");
            }
        } catch (IOException e) {
            logger.error("Failed to read transcript file", e);
            throw new YouTubeException("Failed to read transcript file", e);
        }

        // Cleanup file
        if (!transcriptFile.delete()) {
            logger.error("Cannot remove transcript file: {}", transcriptFile.getAbsolutePath());
            throw new YouTubeException("Cannot remove transcript file");
        }

        VideoTranscript videoTranscript = new VideoTranscript("youtube", videoId, languageCode, LocalDateTime.now(), cleanTranscript(transcript.toString()));
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
        List<String> command = new ArrayList<>();
        command.add("yt-dlp");

        if (StringUtils.hasText(ytDlpProxy)) {
            command.add("--proxy");
            command.add(ytDlpProxy);
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
