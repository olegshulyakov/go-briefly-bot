// Package youtube provides utilities for working with YouTube URLs.
// It includes functionality to validate YouTube URLs, extract video IDs,
// and find all YouTube URLs in a given text.
package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/utils"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/youtube/ytdlp"
)

// extension defines the subtitle format used for transcript conversion.
const extension = "srt"

// ytRegexp compiles the regular expression for matching YouTube URLs.
// Uses the youtube URL Pattern constant and panics if the pattern is invalid.
var ytRegex = regexp.MustCompile(`(?:https?://)?(?:www\.)?(?:youtube\.com/watch\?.*?v=|youtu\.be/)([a-zA-Z0-9_-]{11})`)

// IsValidURL checks if the given text string contains a valid YouTube URL.
// Returns true if the text matches the YouTube URL pattern, false otherwise.
// An empty string will always return false.
func IsValidURL(text string) bool {
	return ytRegex.MatchString(text)
}

// GetID extracts the YouTube video ID from a valid YouTube URL.
// Returns the 11-character video ID if found, or an error if:
// - The input is empty
// - No valid YouTube URL is present
// - No video ID can be extracted from the URL.
func GetID(text string) (string, error) {
	if text == "" {
		return "", errors.New("empty input")
	}
	matches := ytRegex.FindStringSubmatch(text)
	if len(matches) == 0 {
		return "", errors.New("no valid YouTube URL found")
	}
	return matches[1], nil
}

// ExtractURLs finds all valid YouTube URLs in the given text.
// Returns a slice of all matching YouTube URL strings, or an empty slice if no valid URLs are found.
func ExtractURLs(text string) []string {
	return ytRegex.FindAllString(text, -1)
}

// DataLoader represents loader for a Youtube video.
type DataLoader struct {
	url     string
	id      string
	isValid bool

	info       *video.Info
	transcript *video.Transcript
}

// New creates a new DataLoader instance for the given YouTube URL.
// The URL is automatically validated during initialization.
func New(url string) (*DataLoader, error) {
	isValid := IsValidURL(url)
	if !isValid {
		return nil, fmt.Errorf("no valid URL found: %s", url)
	}

	id, err := GetID(url)
	if err != nil {
		return nil, err
	}

	return &DataLoader{url: fmt.Sprintf("https://www.youtube.com/watch?v=%s", id), isValid: true, id: id}, nil
}

// VideoInfo returns the video metadata after successful loading.
// Returns nil if Load() hasn't been called or failed.
func (loader *DataLoader) VideoInfo() *video.Info {
	return loader.info
}

// Transcript returns the processed transcript data after successful loading.
// Returns nil if Load() hasn't been called, transcript loading failed, or no transcript was available.
func (loader *DataLoader) Transcript() *video.Transcript {
	return loader.transcript
}

// Load fetches and processes video data from YouTube.
// Returns error if:
// - URL is invalid
// - Video metadata retrieval fails
// - Transcript download/processing fails
// - File operations encounter errors
//
// The method performs the following steps:
// 1. Validates the URL
// 2. Extracts the video ID
// 3. Fetches video metadata using youtube-dl
// 4. Downloads and processes subtitles (SRT format)
// 5. Cleans and stores the transcript data.
func (loader *DataLoader) Load() error {
	slog.Debug("VideoData download", "url", loader.url)
	defer slog.Debug("VideoData downloaded", "url", loader.url)

	if !loader.isValid {
		return fmt.Errorf("no valid URL found: %s", loader.url)
	}

	slog.Debug("VideoInfo download", "url", loader.url)

	var (
		execOutput []byte
		err        error
	)
	if execOutput, err = ytdlp.Exec([]string{"--dump-json"}, loader.url); err != nil {
		return fmt.Errorf("failed to extract video info: %w", err)
	}

	if err = json.Unmarshal(execOutput, &loader.info); err != nil {
		return fmt.Errorf("failed to parse video info: %w", err)
	}

	slog.Debug("VideoInfo downloaded", "url", loader.url)
	slog.Debug("Transcript load", "url", loader.url)

	execOutput, err = ytdlp.Exec(
		[]string{
			"--no-progress",
			"--skip-download",
			"--write-subs",
			"--write-auto-subs",
			"--convert-subs", extension,
			"--sub-lang", fmt.Sprintf("%s,%s_auto,-live_chat", loader.info.Language, loader.info.Language),
			"--output", filepath.Join(os.TempDir(), fmt.Sprintf("subtitles_%s.%%(ext)s", loader.id)),
		},
		loader.url,
	)
	if err != nil {
		return fmt.Errorf("failed to download transcript: %w\n%s", err, execOutput)
	}
	slog.Debug("Transcript downloaded", "url", loader.url)

	// Read the transcript file
	filename := filepath.Join(os.TempDir(), fmt.Sprintf("subtitles_%s.%s.%s", loader.id, loader.info.Language, extension))
	text, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("no subtitles found: %w", err)
	}
	_ = os.Remove(filename)

	transcript := string(text)
	if transcript, err = utils.CleanSRT(transcript); err != nil {
		return fmt.Errorf("failed to clean transcript file: %w", err)
	}

	loader.transcript = &video.Transcript{
		ID:         loader.info.ID,
		Uploader:   loader.info.Uploader,
		Language:   loader.info.Language,
		Title:      loader.info.Title,
		Thumbnail:  loader.info.Thumbnail,
		Transcript: transcript,
	}

	return nil
}
