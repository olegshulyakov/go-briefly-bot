// Package youtube provides utilities for working with YouTube URLs.
// It includes functionality to validate YouTube URLs, extract video IDs,
// and find all YouTube URLs in a given text.
package youtube

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video"
)

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

// New creates a new DataLoader instance for the given YouTube URL.
// The URL is automatically validated during initialization.
func New(url string) (*video.DataLoader, error) {
	isValid := IsValidURL(url)
	if !isValid {
		return nil, fmt.Errorf("no valid URL found: %s", url)
	}

	id, err := GetID(url)
	if err != nil {
		return nil, err
	}

	return video.NewVideoDataLoader(fmt.Sprintf("https://www.youtube.com/watch?v=%s", id), id, true), nil
}
