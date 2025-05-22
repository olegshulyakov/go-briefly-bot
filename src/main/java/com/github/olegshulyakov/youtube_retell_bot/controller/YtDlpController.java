package com.github.olegshulyakov.youtube_retell_bot.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.olegshulyakov.youtube_retell_bot.dao.VideoInfoEntityRepository;
import com.github.olegshulyakov.youtube_retell_bot.exception.YtDlpException;
import com.github.olegshulyakov.youtube_retell_bot.model.VideoInfo;
import com.github.olegshulyakov.youtube_retell_bot.model.VideoInfoEntity;
import com.github.olegshulyakov.youtube_retell_bot.service.YoutubeUrlValidator;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.util.StringUtils;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.TimeUnit;

@RestController
@RequestMapping("yt-dlp")
public class YtDlpController implements YoutubeUrlValidator {
    private final Logger logger = LoggerFactory.getLogger(YtDlpController.class);
    @Autowired
    private final VideoInfoEntityRepository videoInfoEntityRepository;
    private final ObjectMapper objectMapper;
    private final String ytDlpProxy;

    public YtDlpController(VideoInfoEntityRepository videoInfoEntityRepository, @Value("${YT_DLP_PROXY}") String ytDlpProxy) {
        this.videoInfoEntityRepository = videoInfoEntityRepository;
        this.ytDlpProxy = ytDlpProxy;
        this.objectMapper = new ObjectMapper();
    }

    @GetMapping("/video-info")
    public VideoInfo getVideoInfo(@RequestParam String url) {
        logger.info("Get video info: {}", url);

        Optional<String> id = getYoutubeId(url);
        if (id.isEmpty()) {
            return null; //TODO
        }

        String uuid = VideoInfoEntity.getUuid("youtube", id.get());
        if (videoInfoEntityRepository.existsByUuid(uuid)) {
            logger.debug("Using video info from cache: {}", uuid);
            VideoInfoEntity videoInfoEntity = videoInfoEntityRepository.findByUuid(uuid);
            return videoInfoEntity.toVideoInfo();
        }

        List<String> command = new ArrayList<>();
        command.add("yt-dlp");
        if (StringUtils.hasText(ytDlpProxy)) {
            command.add("--proxy");
            command.add(ytDlpProxy);
        }
        command.add("--dump-json");
        command.add(url);

        logger.debug("Executing command: {}", String.join(" ", command));

        try {
            ProcessBuilder processBuilder = new ProcessBuilder(command);
            Process process = processBuilder.start();

            StringBuilder output = new StringBuilder();
            StringBuilder errorOutput = new StringBuilder();

            try (BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()));
                 BufferedReader errorReader = new BufferedReader(new InputStreamReader(process.getErrorStream()))) {

                String line;
                while ((line = reader.readLine()) != null) {
                    output.append(line).append("\n"); // Add newline for clarity
                }

                String errorLine;
                while ((errorLine = errorReader.readLine()) != null) {
                    errorOutput.append(errorLine).append("\n");
                }
            }

            process.waitFor(30, TimeUnit.SECONDS);
            int exitCode = process.exitValue();
            if (exitCode != 0) {
                logger.warn("yt-dlp finished with exit code {}\n{}", exitCode, errorOutput);
                throw new YtDlpException("yt-dlp command failed with exit code " + exitCode);
            }

            logger.info("Video info: {}", output);

            @SuppressWarnings("unchecked")
            Map<String, Object> jsonMap = objectMapper.readValue(output.toString(), Map.class);
            VideoInfo videoInfo = new VideoInfo(
                    "youtube",
                    (String) jsonMap.get("id"),
                    (String) jsonMap.get("uploader"),
                    (String) jsonMap.get("title"),
                    (String) jsonMap.get("thumbnail")
            );

            videoInfoEntityRepository.save(new VideoInfoEntity(videoInfo));

            return videoInfo;
        } catch (IOException | InterruptedException e) {
            logger.error("Failed to read video info", e);
            throw new YtDlpException("Failed to read video info", e);
        }
    }
}
