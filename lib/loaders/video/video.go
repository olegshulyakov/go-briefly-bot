// Package video provides functionality to retrieve video transcripts
// from videos. It extracts video metadata (title, uploader, thumbnail)
// along with the transcript text in the video's primary language.
package video

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/transcripts"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/ytdlp"
)

// extension defines the subtitle format used for transcript conversion.
const extension = "srt"

// Info represents metadata about a video.
type Info struct {
	ID          string                `json:"id"`           // Unique identifier
	Language    string                `json:"language"`     // Language code
	Uploader    string                `json:"uploader"`     // Name of uploader/channel
	Title       string                `json:"title"`        // Title of the video
	Thumbnail   string                `json:"thumbnail"`    // Thumbnail URL
	Description string                `json:"description"`  // Description of the video
	Duration    int                   `json:"duration"`     // Duration
	Timestamp   int64                 `json:"timestamp"`    // Creation date
	IsLive      bool                  `json:"is_live"`      // Is video live one
	OriginalURL string                `json:"original_url"` // Original URL
	Extractor   string                `json:"extractor"`    // Extractor
	Subtitles   map[string][]Subtitle `json:"subtitles"`    // Subtitles map
}

type Subtitle struct {
	Ext string `json:"ext"` // File extension, e.g. 'vtt', 'srt', 'json3', 'ttml'
	URL string `json:"url"` // File download URL
}

// Transcript represents the complete transcript information for a video,
// including metadata and the transcript text itself.
type Transcript struct {
	ID         string `json:"id"`         // Unique video ID
	Language   string `json:"language"`   // Primary language code of the video (e.g., "en")
	Uploader   string `json:"uploader"`   // Name of the video uploader/channel
	Title      string `json:"title"`      // Title of the video
	Thumbnail  string `json:"thumbnail"`  // URL to the video thumbnail image
	Transcript string `json:"transcript"` // Full text transcript of the video
}

// InfoLoader represents loader interface for a video metadata.
type InfoLoader interface {
	// VideoInfo returns metadata for the video.
	// Returns nil if metadata could not be retrieved.
	VideoInfo() *Info
}

// TranscriptLoader represents loader interface for a video transcript.
type TranscriptLoader interface {
	// Transcript returns the transcript text for the video.
	// Returns nil if the transcript could not be retrieved.
	Transcript() *Transcript
}

// DataLoader represents loader for a video.
type DataLoader struct {
	url     string
	id      string
	isValid bool

	info       *Info
	transcript *Transcript
}

func NewVideoDataLoader(url string, id string, isValid bool) *DataLoader {
	return &DataLoader{url: url, isValid: isValid, id: id}
}

// VideoInfo returns the video metadata after successful loading.
// Returns nil if Load() hasn't been called or failed.
func (loader *DataLoader) VideoInfo() *Info {
	return loader.info
}

// Transcript returns the processed transcript data after successful loading.
// Returns nil if Load() hasn't been called, transcript loading failed, or no transcript was available.
func (loader *DataLoader) Transcript() *Transcript {
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

	lang := loader.info.Language
	if lang == "" {
		subtitlesMap := loader.info.Subtitles
		if len(subtitlesMap) == 0 {
			return errors.New("no subtitles available")
		}

		// Fallback to English if exists
		if _, exists := subtitlesMap["en"]; exists {
			lang = "en"
		} else {
			// Get first available language
			for subtitleLang := range subtitlesMap {
				lang = subtitleLang
				break
			}
		}
	}

	execOutput, err = ytdlp.Exec(
		[]string{
			"--no-progress",
			"--skip-download",
			"--write-subs",
			"--write-auto-subs",
			"--convert-subs", extension,
			"--sub-lang", fmt.Sprintf("%s,%s_auto,-live_chat", lang, lang),
			"--output", filepath.Join(os.TempDir(), fmt.Sprintf("subtitles_%s.%%(ext)s", loader.id)),
		},
		loader.url,
	)
	if err != nil {
		return fmt.Errorf("failed to download transcript: %w\n%s", err, execOutput)
	}
	slog.Debug("Transcript downloaded", "url", loader.url)

	// Read the transcript file
	filename := filepath.Join(os.TempDir(), fmt.Sprintf("subtitles_%s.%s.%s", loader.id, lang, extension))
	text, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("no subtitles found: %w", err)
	}
	_ = os.Remove(filename)

	transcript := string(text)
	if transcript, err = transcripts.CleanSRT(transcript); err != nil {
		return fmt.Errorf("failed to clean transcript file: %w", err)
	}

	loader.transcript = &Transcript{
		ID:         loader.info.ID,
		Uploader:   loader.info.Uploader,
		Language:   lang,
		Title:      loader.info.Title,
		Thumbnail:  loader.info.Thumbnail,
		Transcript: transcript,
	}

	return nil
}
