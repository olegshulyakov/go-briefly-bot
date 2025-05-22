package com.github.olegshulyakov.youtube_retell_bot.service;

import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;

class YoutubeUrlValidatorTest {

    YoutubeUrlValidator validator = new YoutubeUrlValidator() {
        @Override
        public boolean isValidYoutubeUrl(String url) {
            return YoutubeUrlValidator.super.isValidYoutubeUrl(url);
        }
    };

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
    void isValidYoutubeUrl(String url, boolean expected) {
        assertEquals(expected, validator.isValidYoutubeUrl(url), "Test failed for URL: " + url);
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
    void getYoutubeId(String url, String expected) {
        Optional<String> id = validator.getYoutubeId(url);
        if (id.isEmpty()) {
            assertNull(expected, "Test failed for URL: " + url);
        } else {
            assertEquals(expected, id.get(), "Test failed for URL: " + url);
        }
    }
}