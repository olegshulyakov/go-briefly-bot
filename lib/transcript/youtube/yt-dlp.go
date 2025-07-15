// Package youtube provides functionality for interacting with YouTube videos,
// including retrieving video information and transcripts.
//
// The package uses yt-dlp (https://github.com/yt-dlp/yt-dlp) as a backend
// for fetching data from YouTube. Additional yt-dlp options can be configured
// through the YT_DLP_ADDITIONAL_OPTIONS environment variable.
package youtube

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/olegshulyakov/go-briefly-bot/lib/transcript/utils"
)

var (
	// ytDlpAdditionalOptions contains additional command-line options for yt-dlp
	// parsed from YT_DLP_ADDITIONAL_OPTIONS environment variable.
	ytDlpAdditionalOptions []string
)

func init() {
	ytDlpAdditionalOptions = strings.Fields(os.Getenv("YT_DLP_ADDITIONAL_OPTIONS"))
}

// Info represents metadata about a YouTube video.
type Info struct {
	ID        string `json:"id"`        // The unique identifier of the video.
	Language  string `json:"language"`  // The video language.
	Uploader  string `json:"uploader"`  // The name of the video uploader.
	Title     string `json:"title"`     // The title of the video.
	Thumbnail string `json:"thumbnail"` // The URL of the video's thumbnail.
}

// VideoInfo retrieves metadata about a YouTube video from its URL.
// It returns an Info struct containing video details or an error if the operation fails.
// The function validates the URL before attempting to fetch information.
func VideoInfo(videoURL string) (*Info, error) {
	slog.Debug("VideoInfo download", "url", videoURL)
	defer slog.Debug("VideoInfo downloaded", "url", videoURL)

	if !IsValidURL(videoURL) {
		return nil, fmt.Errorf("no valid URL found: %s", videoURL)
	}

	var (
		jsonData  []byte
		videoInfo *Info
		err       error
	)
	if jsonData, err = execYtDlp([]string{"--dump-json"}, videoURL); err != nil {
		return nil, fmt.Errorf("failed to extract video info: %w", err)
	}

	if err = json.Unmarshal(jsonData, &videoInfo); err != nil {
		return nil, fmt.Errorf("failed to parse video info: %w", err)
	}

	return videoInfo, nil
}

// Transcript retrieves the transcript/subtitles for a YouTube video.
// It accepts a video URL and optional language code (defaults to English if empty).
// Returns the cleaned transcript text in SRT format or an error if the operation fails.
// The function automatically removes the downloaded subtitle file after reading it.
func Transcript(videoURL string, languageCode string) (string, error) {
	slog.Debug("Transcript load", "url", videoURL)
	defer slog.Debug("Transcript loaded", "url", videoURL)

	const extension = "srt"
	var (
		videoID    string
		transcript string
		err        error
	)

	if !IsValidURL(videoURL) {
		return "", fmt.Errorf("no valid URL found: %s", videoURL)
	}
	if languageCode != "" {
		languageCode = "en"
	}

	if videoID, err = GetID(videoURL); err != nil {
		return "", err
	}

	output, err := execYtDlp(
		[]string{
			"--no-progress",
			"--skip-download",
			"--write-subs",
			"--write-auto-subs",
			"--convert-subs", extension,
			"--sub-lang", fmt.Sprintf("%s,%s_auto,-live_chat", languageCode, languageCode),
			"--output", filepath.Join(os.TempDir(), fmt.Sprintf("subtitles_%s.%%(ext)s", videoID)),
		},
		videoURL,
	)
	if err != nil {
		return "", fmt.Errorf("failed to download transcript: %w\n%s", err, output)
	}
	slog.Debug("Transcript downloaded", "url", videoURL)

	// Read the transcript file
	filename := filepath.Join(os.TempDir(), fmt.Sprintf("subtitles_%s.%s.%s", videoID, languageCode, extension))
	text, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("no subtitles found: %w", err)
	}
	_ = os.Remove(filename)

	transcript = string(text)
	if transcript, err = utils.CleanSRT(transcript); err != nil {
		return "", fmt.Errorf("failed to clean transcript file: %w", err)
	}

	return transcript, nil
}

// execYtDlp executes the yt-dlp command with the provided arguments and URL.
// It automatically retries the command up to 3 times on failure and includes
// any additional options specified in YT_DLP_ADDITIONAL_OPTIONS.
// Returns the command output or an error if all attempts fail.
func execYtDlp(arguments []string, url string) ([]byte, error) {
	const maxAttempts = 3
	var (
		err    error
		output []byte
		args   = make([]string, 0, len(arguments)+len(ytDlpAdditionalOptions)+1)
	)

	if len(ytDlpAdditionalOptions) > 0 {
		args = append(args, ytDlpAdditionalOptions...)
	}
	args = append(args, arguments...)
	args = append(args, url)

	// Execute with retry
	slog.Debug("Executing yt-dlp", "args", args)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if output, err = exec.Command("yt-dlp", args...).Output(); err == nil {
			break
		}
	}

	if err != nil {
		return output, fmt.Errorf("failed to extract transcript: %w\n%s", err, output)
	}

	return output, nil
}
