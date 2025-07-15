// Package youtube provides utilities for working with YouTube URLs.
// It includes functionality to validate YouTube URLs, extract video IDs,
// and find all YouTube URLs in a given text.
package youtube

import (
	"errors"
	"regexp"
)

// youtubeURLPattern is the regular expression pattern used to match YouTube URLs.
// It supports various formats including:
// - https://www.youtube.com/watch?v=ID,
// - http://youtu.be/ID,
// - www.youtube.com/watch?v=ID,
// - youtube.com/watch?v=ID.
const youtubeURLPattern = `(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]{11})`

// IsValidURL checks if the given text string contains a valid YouTube URL.
// Returns true if the text matches the YouTube URL pattern, false otherwise.
// An empty string will always return false.
func IsValidURL(text string) bool {
	if text == "" {
		return false
	}
	return ytRegexp().MatchString(text)
}

// GetID extracts the YouTube video ID from a valid YouTube URL.
// Returns the 11-character video ID if found, or an error if:
// - The input is not a valid YouTube URL.
// - No video ID could be extracted from the URL.
func GetID(text string) (string, error) {
	if !IsValidURL(text) {
		return "", errors.New("no valid URL found")
	}
	matches := ytRegexp().FindStringSubmatch(text)
	if len(matches) == 0 {
		return "", errors.New("YouTube ID not found")
	}
	return matches[1], nil
}

// ExtractURLs finds all valid YouTube URLs in the given text.
// Returns a slice of all matching YouTube URL strings, or an error if no valid URLs are found.
func ExtractURLs(text string) ([]string, error) {
	if !IsValidURL(text) {
		return nil, errors.New("no valid URL found")
	}
	return ytRegexp().FindAllString(text, -1), nil
}

// ytRegexp compiles and returns the regular expression for matching YouTube URLs.
// Uses the youtubeURLPattern constant and panics if the pattern is invalid.
func ytRegexp() *regexp.Regexp {
	return regexp.MustCompile(youtubeURLPattern)
}
