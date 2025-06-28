package com.github.youtubebriefly.file;

import com.github.youtubebriefly.exception.YouTubeException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
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

    /**
     * Reads SRT file from the current directory.
     * It deletes file after execution.
     * @param fileName file name to read.
     * @return transcript from the file.
     * @throws IOException if fail to read the file.
     */
    public static String readAndDelete(String fileName) throws IOException {
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
        }

        // Cleanup file
        if (!transcriptFile.delete()) {
            logger.error("Cannot remove transcript file: {}", transcriptFile.getAbsolutePath());
            throw new IOException("Cannot remove transcript file");
        }

        return transcript.toString();
    }
}
