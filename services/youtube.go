package services

import (
	"fmt"
	"os"
	"os/exec"
)

func GetTranscript(videoURL string) (string, error) {
	// Use yt-dlp to extract transcript
	cmd := exec.Command("yt-dlp", "--skip-download", "--write-auto-sub", "--sub-format", "srt", "--output", "transcript", videoURL)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to extract transcript: %v", err)
	}

	// Read the transcript file
	transcript, err := os.ReadFile("transcript.en.srt")
	if err != nil {
		return "", fmt.Errorf("failed to read transcript file: %v", err)
	}

	return string(transcript), nil
}