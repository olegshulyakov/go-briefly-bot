package com.github.youtubebriefly.file;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import static org.junit.jupiter.api.Assertions.assertEquals;

public class TranscriptsTest {

    @ParameterizedTest
    @CsvSource({
            "'00:00:00,000 --> 00:00:05,000\nSome text',Some text", // Timeline and text
    })
    void testCleanSRT_TimelineAndText(String transcript, String expected) {
        assertEquals(expected, Transcripts.cleanSRT(transcript));
    }

    @ParameterizedTest
    @CsvSource({
            "'123\n456\nSome text',Some text", // Numeric lines and text
    })
    void testCleanSRT_NumericLinesAndText(String transcript, String expected) {
        assertEquals(expected, Transcripts.cleanSRT(transcript));
    }


    @ParameterizedTest
    @CsvSource({
            "'Line1\nLine1\nLine2','Line1 Line2'", //Duplicate Lines
            "'  Line1   \nLine1',Line1" //Duplicate Lines with whitespace
    })
    void testCleanSRT_DuplicatesRemoved(String transcript, String expected) {
        assertEquals(expected, Transcripts.cleanSRT(transcript));
    }

    @ParameterizedTest
    @CsvSource({
            "'00:00:00,000 --> 00:00:05,000',''",
            "'00:00:10,000 --> 00:00:15,000\n00:00:20,000 --> 00:00:25,000',''"
    })
    void testCleanSRT_OnlyTimelines(String transcript, String expected) {
        assertEquals(expected, Transcripts.cleanSRT(transcript));
    }

    @ParameterizedTest
    @CsvSource({
            "123,''",
            "'456\n789',''"
    })
    void testCleanSRT_OnlyNumbers(String transcript, String expected) {
        assertEquals(expected, Transcripts.cleanSRT(transcript));
    }

    @ParameterizedTest
    @CsvSource({
            "'  Line 1  \nLine 2   ','Line 1 Line 2'", //Leading/trailing whitespace
    })
    void testCleanSRT_LeadingTrailingWhitespace(String transcript, String expected) {
        assertEquals(expected, Transcripts.cleanSRT(transcript));
    }

    @Test
    void testCleanSRT_EmptyInput() {
        assertEquals("", Transcripts.cleanSRT(""));
    }
}