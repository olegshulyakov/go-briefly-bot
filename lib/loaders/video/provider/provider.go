// Package provider provides utilities for working with URLs.
// It includes functionality to validate URLs against regexp, extract video IDs,
// and find all URLs in a given text.
package provider

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video"
)

var (
	Youtube      = ProviderByRegex{regex: regexp.MustCompile(`(?:https?://)?(?:www\.)?(?:youtube\.com/watch\?.*?v=|youtu\.be/)([a-zA-Z0-9_-]{11})`), canonicalUrl: "https://www.youtube.com/watch?v=%s"}
	YoutubeShort = ProviderByRegex{regex: regexp.MustCompile(`(?:https?://)?(?:www\.)?youtube\.com/shorts/([A-Za-z0-9_-]{11})`), canonicalUrl: "https://www.youtube.com/shorts/%s"}
	VkVideo      = ProviderByRegex{regex: regexp.MustCompile(`(?:https?://)?(?:www\.)?vkvideo\.ru/(video-\d+_\d+)`), canonicalUrl: "https://vkvideo.ru/%s"}
)

// ProviderByRegex represents a video provider with regex validation.
type ProviderByRegex struct {
	regex        *regexp.Regexp
	canonicalUrl string
}

// IsValidURL checks if the given text string contains a valid URL by Regex.
// Returns true if the text matches the Regex pattern, false otherwise.
// An empty string will always return false.
func (provider *ProviderByRegex) IsValidURL(text string) bool {
	return provider.regex.MatchString(text)
}

// GetID extracts the video ID from a valid URL.
// Returns the video ID if found, or an error if:
// - The input is empty
// - No valid URL is present
// - No video ID can be extracted from the URL.
func (provider *ProviderByRegex) GetID(text string) (string, error) {
	if text == "" {
		return "", errors.New("empty input")
	}
	matches := provider.regex.FindStringSubmatch(text)
	if len(matches) == 0 {
		return "", errors.New("no valid YouTube URL found")
	}
	return matches[1], nil
}

// ExtractURLs finds all valid URLs in the given text.
// Returns a slice of all matching URL strings, or an empty slice if no valid URLs are found.
func (provider *ProviderByRegex) ExtractURLs(text string) []string {
	return provider.regex.FindAllString(text, -1)
}

// Builds a new DataLoader instance for the given URL.
// The URL is automatically validated during initialization.
func (provider *ProviderByRegex) BuildDataLoader(url string) (*video.DataLoader, error) {
	isValid := provider.IsValidURL(url)
	if !isValid {
		return nil, fmt.Errorf("no valid URL found: %s", url)
	}

	id, err := provider.GetID(url)
	if err != nil {
		return nil, err
	}

	return video.NewVideoDataLoader(fmt.Sprintf(provider.canonicalUrl, id), id, true), nil
}
