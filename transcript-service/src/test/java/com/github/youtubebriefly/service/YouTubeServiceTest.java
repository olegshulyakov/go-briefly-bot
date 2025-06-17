package com.github.youtubebriefly.service;

import com.github.youtubebriefly.exception.YouTubeException;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;

class YouTubeServiceTest {

    @ParameterizedTest
    @CsvSource({
            "https://www.youtube.com/watch?v=dQw4w9WgXcQ,true",
            "http://www.youtube.com/watch?v=dQw4w9WgXcQ,true",
            "https://youtu.be/dQw4w9WgXcQ,true",
            "http://youtu.be/dQw4w9WgXcQ,true",
            "https://www.youtube.com/watch?v=dQw4w9WgXcQ123,false",
            "https://www.youtube.com/watch?v=dQw4w9WgXcQ&feature=youtu.be,false",
            "https://www.youtube.com/embed/dQw4w9WgXcQ,false",
            "https://www.youtube.com/watch?v=invalid-id,false",
            "https://www.youtube.com/watch?v=dQw4w9WgXcQ&feature=somethingelse,false",
            "https://www.youtube.com/watch?v=,false"
    })
    void isValidUrl(String url, boolean expected) {
        assertEquals(expected, YouTubeService.isValidUrl(url), "Test failed for URL: " + url);
    }

    @ParameterizedTest
    @CsvSource({
            "https://www.youtube.com/watch?v=dQw4w9WgXcQ,dQw4w9WgXcQ",
            "http://www.youtube.com/watch?v=dQw4w9WgXcQ,dQw4w9WgXcQ",
            "https://youtu.be/dQw4w9WgXcQ,dQw4w9WgXcQ",
            "http://youtu.be/dQw4w9WgXcQ,dQw4w9WgXcQ",
            "https://www.youtube.com/watch?v=dQw4w9WgXcQ123,",
            "https://www.youtube.com/watch?v=dQw4w9WgXcQ&feature=youtu.be,",
            "https://www.youtube.com/embed/dQw4w9WgXcQ,",
            "https://www.youtube.com/watch?v=invalid-id,",
            "https://www.youtube.com/watch?v=dQw4w9WgXcQ&feature=somethingelse,",
            "https://www.youtube.com/watch?v=,"
    })
    void getVideoId(String url, String expected) {
        try {
            assertEquals(expected, YouTubeService.getVideoId(url), "Test failed for URL: " + url);
        } catch (YouTubeException e) {
            assertNull(expected, "Test failed for URL: " + url);
        }
    }
}