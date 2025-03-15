package services

import (
	"fmt"
	"os"
	"os/exec"
	"youtube-retell-bot/config"
)

func GetTranscript(videoURL string) (string, error) {
	config.Logger.Debugf("Transcript extract: %v", videoURL)

	cmd := exec.Command("yt-dlp", "--skip-download", "--write-auto-sub", "--convert-subs", "srt", "--sub-langs", "ru,ru_auto,-live_chat", "--output", "transcript", videoURL)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to extract transcript: %v", err)
	}

	config.Logger.Debugf("Transcript downloaded: %v", videoURL)

	// Read the transcript file
	transcript, err := os.ReadFile("transcript.ru.srt")
	if err != nil {
		return "", fmt.Errorf("failed to read transcript file: %v", err)
	}

	config.Logger.Debugf("Transcript extracted: %v", videoURL)

	return string(transcript), nil
}