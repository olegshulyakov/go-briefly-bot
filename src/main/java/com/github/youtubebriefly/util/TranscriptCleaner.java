package com.github.youtubebriefly.util;

import java.util.HashSet;
import java.util.Set;
import java.util.regex.Pattern;

public interface TranscriptCleaner {
    Pattern NUMERIC_LINES_REGEX = Pattern.compile("^\\d+$");
    Pattern TIMELINE_REGEX = Pattern.compile("^\\d{2}:\\d{2}:\\d{2},\\d{3} --> \\d{2}:\\d{2}:\\d{2},\\d{3}$");

    default String cleanTranscript(String transcript) {
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
