package com.github.youtubebriefly;

import com.github.youtubebriefly.file.TranscriptFiles;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Application implements CommandLineRunner {

    /**
     * Main class for the application.
     *
     * @param args Command line arguments
     */
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public void run(String... args) {
        TranscriptFiles.cleanUpOldFiles();
    }
}
