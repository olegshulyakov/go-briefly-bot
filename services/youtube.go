// Package services provides functionality for interacting with YouTube videos,
// such as extracting video information and transcripts.
package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"youtube-retell-bot/config"
)

// VideoInfo represents metadata about a YouTube video.
type VideoInfo struct {
	Id        string `json:"id"`        // The unique identifier of the video.
	Uploader  string `json:"uploader"`  // The name of the video uploader.
	Title     string `json:"title"`     // The title of the video.
	Thumbnail string `json:"thumbnail"` // The URL of the video's thumbnail.
}

// GetYoutubeVideoInfo retrieves metadata about a YouTube video using its URL.
//
// The function uses the `yt-dlp` command-line tool to extract video information
// in JSON format and then parses it into a VideoInfo struct.
//
// Parameters:
//   - videoURL: The URL of the YouTube video.
//
// Returns:
//   - A VideoInfo struct containing the video's metadata.
//   - An error if the video information cannot be extracted or parsed.
//
// Example:
//
//	videoInfo, err := GetYoutubeVideoInfo("https://www.youtube.com/watch?v=example")
//	if err != nil {
//	    log.Errorf("Failed to get video info: %v", err)
//	}
//	fmt.Printf("Video Title: %s\n", videoInfo.Title)
//
// Notes:
//   - The function relies on the `yt-dlp` tool being installed and accessible in the system's PATH.
//   - Logging is performed using the `config.Logger` for debugging and error tracking.
func GetYoutubeVideoInfo(videoURL string) (VideoInfo, error) {
	var videoInfo VideoInfo

	config.Logger.Debugf("VideoInfo download: %v", videoURL)

	cmd := exec.Command("yt-dlp", "--dump-json", videoURL)

	output, err := cmd.Output()
	if err != nil {
		return videoInfo, fmt.Errorf("failed to extract video info: %v\n", err)
	}

	err = json.Unmarshal(output, &videoInfo)
	if err != nil {
		return videoInfo, fmt.Errorf("failed to parse video info: %v\n", err)
	}

	config.Logger.Debugf("VideoInfo downloaded: %v", videoURL)

	return videoInfo, nil
}

// GetYoutubeTranscript retrieves the transcript of a YouTube video using its URL.
//
// The function uses the `yt-dlp` command-line tool to extract the transcript
// in SRT format, reads the transcript file, and then deletes the file.
//
// Parameters:
//   - videoURL: The URL of the YouTube video.
//
// Returns:
//   - A string containing the transcript of the video.
//   - An error if the transcript cannot be extracted, read, or the file cannot be deleted.
//
// Example:
//
//	transcript, err := GetYoutubeTranscript("https://www.youtube.com/watch?v=example")
//	if err != nil {
//	    log.Errorf("Failed to get transcript: %v", err)
//	}
//	fmt.Println("Transcript:", transcript)
//
// Notes:
//   - The function relies on the `yt-dlp` tool being installed and accessible in the system's PATH.
//   - The transcript is extracted in Russian (`ru` and `ru_auto`) and saved as an SRT file.
//   - The transcript file is deleted after reading to clean up temporary files.
//   - Logging is performed using the `config.Logger` for debugging and error tracking.
func GetYoutubeTranscript(videoURL string) (string, error) {
	config.Logger.Debugf("Transcript extract: %v", videoURL)

	cmd := exec.Command("yt-dlp", "--skip-download", "--write-auto-sub", "--convert-subs", "srt", "--sub-langs", "ru,ru_auto,-live_chat", "--output", "transcript", videoURL)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to extract transcript: %v", err)
	}

	// Read the transcript file
	transcript, err := os.ReadFile("transcript.ru.srt")
	if err != nil {
		return "", fmt.Errorf("failed to read transcript file: %v", err)
	}

	err = os.Remove("transcript.ru.srt")
	if err != nil {
		return "", fmt.Errorf("failed to delete transcript file: %v", err)
	}

	config.Logger.Debugf("Transcript extracted: %v", videoURL)

	return string(transcript), nil
}
