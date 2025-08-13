package utils

import (
	"regexp"
)

// YouTubeVideoIDRegex is the compiled regex for extracting YouTube video IDs.
// It matches the pattern provided in the system design.
var YouTubeVideoIDRegex = regexp.MustCompile(`(?i)(?:https?://)?(?:www\.)?(?:youtube\.com/watch\?.*?v=|youtu\.be/)([a-zA-Z0-9_-]{11})`)

// ExtractYouTubeVideoIDs extracts YouTube video IDs from a given text.
func ExtractYouTubeVideoIDs(text string) []string {
	matches := YouTubeVideoIDRegex.FindAllStringSubmatch(text, -1)
	var ids []string
	seen := make(map[string]bool) // To avoid duplicates

	for _, match := range matches {
		if len(match) > 1 {
			id := match[1]
			// Basic validation: ensure it's exactly 11 characters and not already added
			if len(id) == 11 && !seen[id] {
				ids = append(ids, id)
				seen[id] = true
			}
		}
	}
	return ids
}

// IsValidYouTubeVideoID checks if a string is a valid YouTube video ID format.
func IsValidYouTubeVideoID(id string) bool {
	// Basic check: length and character set
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]{11}$`, id)
	return matched
}

// ConstructYouTubeURL constructs a full YouTube URL from a video ID.
func ConstructYouTubeURL(videoID string) string {
	if !IsValidYouTubeVideoID(videoID) {
		return "" // Return empty string for invalid ID
	}
	return "https://www.youtube.com/watch?v=" + videoID
}

// ExtractYouTubeURLs is a convenience function that extracts full URLs.
func ExtractYouTubeURLs(text string) []string {
	ids := ExtractYouTubeVideoIDs(text)
	var urls []string
	for _, id := range ids {
		url := ConstructYouTubeURL(id)
		if url != "" {
			urls = append(urls, url)
		}
	}
	return urls
}