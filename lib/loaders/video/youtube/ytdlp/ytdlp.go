// Package ytdlp provides functionality for interacting with YouTube videos,
// including retrieving video information and transcripts.
//
// The package uses yt-dlp (https://github.com/yt-dlp/yt-dlp) as a backend
// for fetching data from YouTube. Additional yt-dlp options can be configured
// through the YT_DLP_ADDITIONAL_OPTIONS environment variable.
package ytdlp

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/utils"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/youtube"
)

// DataLoader represents loader for a Youtube video.
type DataLoader struct{}

func New() *DataLoader {
	return &DataLoader{}
}

// Transcript retrieves and returns the complete transcript information
// for a YouTube video given its URL.
//
// Parameters:
//   - videoURL: The full YouTube video URL (e.g., "https://youtube.com/watch?v=...")
//
// Returns:
//   - *VideoTranscript containing all video metadata and transcript text
//   - error if any step of the process fails (video info or transcript retrieval)
//
// Example:
//
//	transcript, err := Transcript("https://youtube.com/watch?v=dQw4w9WgXcQ")
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(transcript.Title)       // Print video title
//	fmt.Println(transcript.Transcript)  // Print transcript text
func (loader *DataLoader) Transcript(videoURL string) (*video.Transcript, error) {
	videoInfo, err := loader.VideoInfo(videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get video info: %s, %w", videoURL, err)
	}

	transcript, err := loader.transcript(videoURL, videoInfo.Language)
	if err != nil {
		return nil, fmt.Errorf("failed to get video transcript: %s, %w", videoURL, err)
	}

	videoTranscript := &video.Transcript{
		ID:         videoInfo.ID,
		Uploader:   videoInfo.Uploader,
		Language:   videoInfo.Language,
		Title:      videoInfo.Title,
		Thumbnail:  videoInfo.Thumbnail,
		Transcript: transcript,
	}

	return videoTranscript, nil
}

// VideoInfo retrieves metadata about a YouTube video from its URL.
// It returns an Info struct containing video details or an error if the operation fails.
// The function validates the URL before attempting to fetch information.
func (loader *DataLoader) VideoInfo(videoURL string) (*video.Info, error) {
	slog.Debug("VideoInfo download", "url", videoURL)
	defer slog.Debug("VideoInfo downloaded", "url", videoURL)

	if err := loader.isValid(videoURL); err != nil {
		return nil, err
	}

	var (
		jsonData  []byte
		videoInfo *video.Info
		err       error
	)
	if jsonData, err = loader.exec([]string{"--dump-json"}, videoURL); err != nil {
		return nil, fmt.Errorf("failed to extract video info: %w", err)
	}

	if err = json.Unmarshal(jsonData, &videoInfo); err != nil {
		return nil, fmt.Errorf("failed to parse video info: %w", err)
	}

	return videoInfo, nil
}

// transcript retrieves the transcript/subtitles for a YouTube video.
// It accepts a video URL and optional language code (defaults to English if empty).
// Returns the cleaned transcript text in SRT format or an error if the operation fails.
// The function automatically removes the downloaded subtitle file after reading it.
func (loader *DataLoader) transcript(videoURL string, languageCode string) (string, error) {
	slog.Debug("Transcript load", "url", videoURL)
	defer slog.Debug("Transcript loaded", "url", videoURL)

	const extension = "srt"

	if err := loader.isValid(videoURL); err != nil {
		return "", err
	}
	if languageCode != "" {
		languageCode = "en"
	}

	videoID, err := youtube.GetID(videoURL)
	if err != nil {
		return "", err
	}

	output, err := loader.exec(
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

	transcript := string(text)
	if transcript, err = utils.CleanSRT(transcript); err != nil {
		return "", fmt.Errorf("failed to clean transcript file: %w", err)
	}

	return transcript, nil
}

// exec executes the yt-dlp command with the provided arguments and URL.
// It automatically retries the command up to 3 times on failure and includes
// any additional options specified in YT_DLP_ADDITIONAL_OPTIONS.
// Returns the command output or an error if all attempts fail.
func (loader *DataLoader) exec(arguments []string, url string) ([]byte, error) {
	const maxAttempts = 3
	// ytDlpAdditionalOptions contains additional command-line options for yt-dlp
	// parsed from YT_DLP_ADDITIONAL_OPTIONS environment variable.
	var ytDlpAdditionalOptions = strings.Fields(os.Getenv("YT_DLP_ADDITIONAL_OPTIONS"))
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

func (loader *DataLoader) isValid(videoURL string) error {
	if !youtube.IsValidURL(videoURL) {
		return fmt.Errorf("no valid URL found: %s", videoURL)
	}
	return nil
}
