// Package services provides functionality for interacting with YouTube videos,
// such as extracting video information and transcripts.
package services

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"youtube-retell-bot/config"
)

// Define the regex pattern for YouTube URLs
var YoutubeUrlPattern = `(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]{11})`

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

	args := []string{
		"--dump-json",
	}

	if config.Configuration.YtDlpProxy != "" {
		args = append(args[:], "--proxy", config.Configuration.YtDlpProxy)
	}

	if videoURL != "" {
		args = append(args[:], videoURL)
	} else {
		return videoInfo, fmt.Errorf("videoURL is empty")
	}

	cmd := exec.Command("yt-dlp", args...)

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
func GetYoutubeTranscript(videoURL string, languageCode string) (string, error) {
	config.Logger.Debugf("Transcript extract: %v", videoURL)

	args := []string{
		"--no-progress",
		"--skip-download",
		"--write-subs",
		"--write-auto-subs",
		"--convert-subs", "srt",
	}

	if config.Configuration.YtDlpProxy != "" {
		args = append(args[:], "--proxy", config.Configuration.YtDlpProxy)
	}

	if languageCode != "" {
		args = append(args[:], "--sub-lang", fmt.Sprintf("%s,%s_auto,-live_chat", languageCode, languageCode))
	}

	if videoURL != "" {
		args = append(args[:], "--output", fmt.Sprintf("subtitles_%s.%%(ext)s", videoURL[len(videoURL)-11:]), videoURL)
	} else {
		return "", fmt.Errorf("videoURL is empty")
	}

	cmd := exec.Command("yt-dlp", args...)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to extract transcript: %v\n%s", err, output)
	}

	// Generate the file name
	fileName := fmt.Sprintf("subtitles_%s.%s.srt", videoURL[len(videoURL)-11:], languageCode)

	// Read the transcript file
	transcript, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("No subtitles to found: %v", err)
	}

	err = os.Remove(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to delete transcript file: %v", err)
	}

	config.Logger.Debugf("Transcript extracted: %v", videoURL)

	return string(transcript), nil
}

// IsValidYouTubeURL checks if the provided text contains a valid YouTube URL.
//
// It compiles a regular expression defined by YoutubeUrlPattern and uses it to
// determine whether the input text contains a YouTube URL.
//
// Parameters:
//   - text: The string to check for a YouTube URL.
//
// Returns:
//   - true if the text contains a valid YouTube URL, false otherwise.
func IsValidYouTubeURL(text string) bool {
	// Compile the regex
	re, err := regexp.Compile(YoutubeUrlPattern)
	if err != nil {
		return false
	}

	// Check if the text contains a YouTube URL
	if !re.MatchString(text) {
		return false
	}
	return true
}

// ExtractAllYouTubeURLs extracts all YouTube URLs from the given text.
// It uses a regular expression to find all matching URLs.
//
// Parameters:
//   - text: The string to extract YouTube URLs from.
//
// Returns:
//   - A slice of strings containing all the YouTube URLs found in the text.
//   - An error if there was an error compiling the regex.
//
// Example:
//
//	urls, err := ExtractAllYouTubeURLs("Check out this video: https://www.youtube.com/watch?v=dQw4w9WgXcQ and another one at https://youtu.be/abcdefg123")
//	if err != nil {
//	    log.Errorf("Error extracting URLs: %v", err)
//	}
//	fmt.Println("URLs:", urls)
func ExtractAllYouTubeURLs(text string) ([]string, error) {
	// Compile the regex
	re, err := regexp.Compile(YoutubeUrlPattern)
	if err != nil {
		return nil, fmt.Errorf("Error compiling regex: %v", err)
	}

	// Check if the text contains a YouTube URL
	if !re.MatchString(text) {
		return nil, fmt.Errorf("No valid URL found")
	}

	// Find all YouTube URLs in text
	return re.FindAllString(text, -1), nil
}
