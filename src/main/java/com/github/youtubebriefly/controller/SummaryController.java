package com.github.youtubebriefly.controller;

import com.github.youtubebriefly.model.SummaryRequest;
import com.github.youtubebriefly.model.SummaryResponse;
import com.github.youtubebriefly.service.SummaryService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/summary")
@RequiredArgsConstructor
public class SummaryController {

    private final SummaryService summaryService;

    @PostMapping("/summarize")
    public ResponseEntity<SummaryResponse> summarize(@RequestBody SummaryRequest request) {
        String summary = summaryService.generateSummary(request.text(), request.languageCode());
        return ResponseEntity.ok(new SummaryResponse(summary, request.languageCode()));
    }
}
