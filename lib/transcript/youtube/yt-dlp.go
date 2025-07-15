package youtube

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/olegshulyakov/go-briefly-bot/lib/transcript/utils"
)

var (
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

// VideoInfo retrieves metadata about a YouTube video using its URL.
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
//	videoInfo, err := VideoInfo("https://www.youtube.com/watch?v=example")
//	if err != nil {
//	    log.Errorf("Failed to get video info: %v", err)
//	}
//	fmt.Printf("Video Title: %s\n", videoInfo.Title)
//
// Notes:
//   - The function relies on the `yt-dlp` tool being installed and accessible in the system's PATH.
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

// Transcript retrieves the transcript of a YouTube video using its URL.
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
//	transcript, err := Transcript("https://www.youtube.com/watch?v=example")
//	if err != nil {
//	    log.Errorf("Failed to get transcript: %v", err)
//	}
//	fmt.Println("Transcript:", transcript)
//
// Notes:
//   - The function relies on the `yt-dlp` tool being installed and accessible in the system's PATH.
//   - The transcript is extracted in Russian (`ru` and `ru_auto`) and saved as an SRT file.
//   - The transcript file is deleted after reading to clean up temporary files.
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
			"--output", fmt.Sprintf("subtitles_%s.%%(ext)s", videoID),
		},
		videoURL,
	)
	if err != nil {
		return "", fmt.Errorf("failed to download transcript: %w\n%s", err, output)
	}

	// Read the transcript file
	if transcript, err = utils.ReadAndRemoveFile(fmt.Sprintf("subtitles_%s.%s.%s", videoID, languageCode, extension)); err != nil {
		return "", err
	}

	slog.Debug("Transcript downloaded", "url", videoURL)

	if transcript, err = utils.CleanSRT(transcript); err != nil {
		return "", fmt.Errorf("failed to clean transcript file: %w", err)
	}

	return transcript, nil
}

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
