package com.github.youtubebriefly.util;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.regex.Pattern;
import java.util.stream.Stream;

public class TranscriptFiles {
    private static final Logger logger = LoggerFactory.getLogger(TranscriptFiles.class);

    private TranscriptFiles() {
    }

    /**
     * Cleans up old SRT files from the current directory.
     * It deletes all files that have an extension matching .srt.
     */
    public static void cleanUpOldFiles() {
        try (Stream<Path> paths = Files.walk(Paths.get("."))) {
            paths.forEach(filePath -> {
                if (!Pattern.matches(".+\\.srt", filePath.getFileName().toString())) {
                    return;
                }

                try {
                    Files.delete(filePath);
                    logger.debug("Deleted: {}", filePath.toAbsolutePath());
                } catch (IOException e) {
                    logger.warn("Error deleting: {}", filePath.toAbsolutePath());
                }
            });
        } catch (IOException e) {
            logger.error("Error during cleanup", e);
        }
    }
}
