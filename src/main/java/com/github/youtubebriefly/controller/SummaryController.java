package com.github.youtubebriefly.controller;

import com.github.youtubebriefly.model.SummaryRequest;
import com.github.youtubebriefly.model.SummaryResponse;
import com.github.youtubebriefly.service.SummaryService;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * REST controller for handling text summarization requests.
 * <p>
 * This controller exposes an endpoint for receiving text and a language code,
 * and returns a summary of the text in the specified language.
 */
@RestController
@RequestMapping("/summary")
@RequiredArgsConstructor
public class SummaryController {

    private static final Logger logger = LoggerFactory.getLogger(SummaryController.class);

    private final SummaryService summaryService;

    /**
     * Generates a summary of the provided text.
     *
     * @param request The {@link SummaryRequest} request, containing the text to summarize and the desired language code.
     * @return A {@link ResponseEntity} containing the {@link SummaryResponse} object with the generated summary and language code,
     * or a {@link ResponseEntity} with an appropriate error status and message if the request is invalid or an error occurs during summarization.
     * @throws IllegalArgumentException if the input `request` is null.
     */
    @PostMapping("/summarize")
    public ResponseEntity<SummaryResponse> summarize(@RequestBody SummaryRequest request) {
        if (request == null) {
            throw new IllegalArgumentException("Request body cannot be null.");
        }

        try {
            String summary = summaryService.generateSummary(request.text(), request.languageCode());
            return ResponseEntity.ok(new SummaryResponse(summary, request.languageCode()));
        } catch (Exception e) {
            logger.warn("Error summarizing text.", e);
            return ResponseEntity.internalServerError().body(new SummaryResponse("Error generating summary.", request.languageCode()));
        }
    }
}
