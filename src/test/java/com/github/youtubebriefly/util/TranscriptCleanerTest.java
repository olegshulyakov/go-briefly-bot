package com.github.youtubebriefly.util;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import static org.junit.jupiter.api.Assertions.assertEquals;

public class TranscriptCleanerTest {
    TranscriptCleaner transcriptCleaner = new TranscriptCleaner() {
        @Override
        public String cleanTranscript(String text) {
            return TranscriptCleaner.super.cleanTranscript(text);
        }
    };

    @ParameterizedTest
    @CsvSource({
            "'00:00:00,000 --> 00:00:05,000\nSome text',Some text", // Timeline and text
    })
    void testCleanTranscript_TimelineAndText(String transcript, String expected) throws Exception {
        assertEquals(expected, transcriptCleaner.cleanTranscript(transcript));
    }

    @ParameterizedTest
    @CsvSource({
            "'123\n456\nSome text',Some text", // Numeric lines and text
    })
    void testCleanTranscript_NumericLinesAndText(String transcript, String expected) throws Exception {
        assertEquals(expected, transcriptCleaner.cleanTranscript(transcript));
    }


    @ParameterizedTest
    @CsvSource({
            "'Line1\nLine1\nLine2','Line1 Line2'", //Duplicate Lines
            "'  Line1   \nLine1',Line1" //Duplicate Lines with whitespace
    })
    void testCleanTranscript_DuplicatesRemoved(String transcript, String expected) throws Exception{
        assertEquals(expected, transcriptCleaner.cleanTranscript(transcript));
    }

    @ParameterizedTest
    @CsvSource({
            "'00:00:00,000 --> 00:00:05,000',''",
            "'00:00:10,000 --> 00:00:15,000\n00:00:20,000 --> 00:00:25,000',''"
    })
    void testCleanTranscript_OnlyTimelines(String transcript, String expected) throws Exception{
        assertEquals(expected, transcriptCleaner.cleanTranscript(transcript));
    }

    @ParameterizedTest
    @CsvSource({
            "123,''",
            "'456\n789',''"
    })
    void testCleanTranscript_OnlyNumbers(String transcript, String expected) throws Exception {
        assertEquals(expected, transcriptCleaner.cleanTranscript(transcript));
    }

    @ParameterizedTest
    @CsvSource({
            "'  Line 1  \nLine 2   ','Line 1 Line 2'", //Leading/trailing whitespace
    })
    void testCleanTranscript_LeadingTrailingWhitespace(String transcript, String expected) throws Exception {
        assertEquals(expected, transcriptCleaner.cleanTranscript(transcript));
    }

    @Test
    void testCleanTranscript_EmptyInput() throws Exception {
        assertEquals("", transcriptCleaner.cleanTranscript(""));
    }
}