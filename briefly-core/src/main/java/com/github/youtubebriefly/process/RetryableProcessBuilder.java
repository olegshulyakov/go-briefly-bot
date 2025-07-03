package com.github.youtubebriefly.process;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.util.Collection;
import java.util.Collections;
import java.util.LinkedList;
import java.util.List;
import java.util.concurrent.TimeUnit;

public class RetryableProcessBuilder {
    private static final Logger logger = LoggerFactory.getLogger(RetryableProcessBuilder.class);

    private final List<String> command = new LinkedList<>();

    public RetryableProcessBuilder(String command) {
        this.command.add(command);
    }

    public RetryableProcessBuilder addOption(String... options) {
        Collections.addAll(this.command, options);
        return this;
    }

    public RetryableProcessBuilder addOption(Collection<? extends String> options) {
        this.command.addAll(options);
        return this;
    }

    public String exec() throws IOException {
        return this.exec(30);
    }

    public String exec(int timeout) throws IOException {
        StringBuilder output = new StringBuilder();
        StringBuilder errorOutput = new StringBuilder();
        Process process = this.start();

        try (BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()));
             BufferedReader errorReader = new BufferedReader(new InputStreamReader(process.getErrorStream()))) {

            String line;
            while ((line = reader.readLine()) != null) {
                output.append(line).append("\n");
            }

            String errorLine;
            while ((errorLine = errorReader.readLine()) != null) {
                errorOutput.append(errorLine).append("\n");
            }
        }

        try {
            process.waitFor(timeout, TimeUnit.SECONDS);
        } catch (InterruptedException e) {
            throw new IOException(e);
        }

        int exitCode = process.exitValue();
        if (exitCode != 0) {
            logger.debug("Command failed with exit code {}\n{}", exitCode, errorOutput);
            throw new IOException(String.format("Command failed with exit code %s\n%s", exitCode, errorOutput));
        }

        return output.toString();
    }

    public String execWithRetry(int maxAttempts) throws IOException {
        if (maxAttempts < 1) {
            throw new IllegalArgumentException("maxAttempts cannot be less than 1");
        }

        int attempt = 0;
        IOException ex = null;

        while (attempt <= maxAttempts + 1) {
            attempt++;

            try {
                return this.exec();
            } catch (Exception e) {
                logger.warn(String.format("Failed : %s, retrying.", String.join(" ", this.command)), e);
                ex = new IOException(e);
            }
        }

        throw ex;
    }

    private Process start() throws IOException {
        logger.debug("Executing command: {}", String.join(" ", this.command));
        return new ProcessBuilder(command).start();
    }
}
