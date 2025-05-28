package com.github.olegshulyakov.youtubebriefly;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.regex.Pattern;
import java.util.stream.Stream;

@SpringBootApplication
public class Application {
    private static final Logger logger = LoggerFactory.getLogger(Application.class);

    /**
     * Main class for the application.
     *
     * @param args Command line arguments
     */
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);

        cleanUpOldFiles();
    }

    /**
     * Cleans up old SRT files from the current directory.
     * It deletes all files that have an extension matching .srt.
     */
    private static void cleanUpOldFiles() {
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
