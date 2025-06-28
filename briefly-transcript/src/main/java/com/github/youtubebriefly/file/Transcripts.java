package com.github.youtubebriefly.file;

import java.util.HashSet;
import java.util.Set;
import java.util.regex.Pattern;

/**
 * Interface for cleaning transcripts by removing unwanted lines such as numeric-only lines,
 * timeline markers, and duplicate entries.
 */
public class Transcripts {

    private Transcripts() {
    }

    /**
     * Regular expression pattern to match lines containing only numbers.
     */
    private static final Pattern NUMERIC_LINES_REGEX = Pattern.compile("^\\d+$");

    /**
     * Regular expression pattern to match lines representing timeline markers
     * (e.g., "00:01:23,456 --> 00:01:25,789").
     */
    private static final Pattern TIMELINE_REGEX = Pattern.compile("^\\d{2}:\\d{2}:\\d{2},\\d{3} --> \\d{2}:\\d{2}:\\d{2},\\d{3}$");

    /**
     * Cleans a transcript string by removing numeric-only lines, timeline markers, and duplicate lines.
     *
     * @param transcript The input transcript string.
     * @return The cleaned transcript string.
     * @throws IllegalArgumentException if the input transcript is null.
     */
    public static String cleanSRT(String transcript) {
        if (transcript == null) {
            throw new IllegalArgumentException("Transcript cannot be null.");
        }

        StringBuilder sb = new StringBuilder();
        Set<String> seen = new HashSet<>();
        String[] lines = transcript.split("\n");

        for (int i = 0; i < lines.length; i++) {
            String line = lines[i].trim();

            // Skip empty lines
            if (line.isEmpty()) {
                continue;
            }

            // Skip timeline lines
            if (TIMELINE_REGEX.matcher(line).matches()) {
                continue;
            }

            // Skip numeric lines
            if (NUMERIC_LINES_REGEX.matcher(line).matches()) {
                continue;
            }

            // Skip duplicate lines
            if (seen.contains(line)) {
                continue;
            }

            // Write the line to output
            seen.add(line);
            sb.append(line);

            // Append space between lines
            if (i < lines.length - 1) {
                sb.append(" ");
            }
        }
        return sb.toString().trim();
    }
}
