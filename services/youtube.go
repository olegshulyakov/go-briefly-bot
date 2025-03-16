package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"youtube-retell-bot/config"
)

type YouTubeResult struct {
	Id         string `json:"id"`
	Uploader   string `json:"uploader"`
	Title      string `json:"title"`
	Thumbnail  string `json:"thumbnail"`
	Transcript string
}

func GetTranscript(videoURL string) (YouTubeResult, error) {
	var youTubeResult YouTubeResult

	config.Logger.Debugf("Video Info extract: %v", videoURL)

	cmd := exec.Command("yt-dlp", "--dump-json", videoURL)

	output, err := cmd.Output()
	if err != nil {
		return youTubeResult, fmt.Errorf("failed to extract video info: %v\n", err)
	}

	err = json.Unmarshal(output, &youTubeResult)
	if err != nil {
		return youTubeResult, fmt.Errorf("failed to parse video info: %v\n", err)
	}

	config.Logger.Debugf("Transcript extract: %v", videoURL)

	cmd = exec.Command("yt-dlp", "--skip-download", "--write-auto-sub", "--convert-subs", "srt", "--sub-langs", "ru,ru_auto,-live_chat", "--output", "transcript", videoURL)
	err = cmd.Run()
	if err != nil {
		return youTubeResult, fmt.Errorf("failed to extract transcript: %v", err)
	}

	config.Logger.Debugf("Transcript downloaded: %v", videoURL)

	// Read the transcript file
	transcript, err := os.ReadFile("transcript.ru.srt")
	if err != nil {
		return youTubeResult, fmt.Errorf("failed to read transcript file: %v", err)
	}
	youTubeResult.Transcript = string(transcript)

	err = os.Remove("transcript.ru.srt")
	if err != nil {
		return youTubeResult, fmt.Errorf("failed to delete transcript file: %v", err)
	}

	config.Logger.Debugf("Transcript extracted: %v", videoURL)

	return youTubeResult, nil
}
