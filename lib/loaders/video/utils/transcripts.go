package utils

import (
	"regexp"
	"strings"
)

// CleanSRT processes SRT (SubRip Subtitle) formatted text and returns a cleaned version.
// It performs the following operations:
//   - Removes empty lines
//   - Removes timeline markers (e.g., "00:00:00,000 --> 00:00:00,000")
//   - Removes subtitle sequence numbers
//   - Deduplicates repeated lines
//   - Joins remaining lines with spaces
//
// Parameters:
//   - text: The input SRT formatted text to be cleaned
//
// Returns:
//   - string: The cleaned text with all lines joined by spaces
//   - error: Non-nil if any error occurs during string building
//
// Example:
//
//	cleaned, err := CleanSRT("1\n00:00:00,000 --> 00:00:02,000\nHello\nHello\n\n2\n00:00:03,000 --> 00:00:05,000\nWorld")
//	// Returns: "Hello World", nil
func CleanSRT(text string) (string, error) {
	timelineRegex := regexp.MustCompile(`^\d{2}:\d{2}:\d{2},\d{3} --> \d{2}:\d{2}:\d{2},\d{3}$`)
	numericLinesRegex := regexp.MustCompile(`^\d+$`)
	specialCharactersAtEndRegex := regexp.MustCompile(`(\\\w)+$`)

	var sb strings.Builder
	seen := make(map[string]bool)
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)

		// Skip timeline lines (00:00:00,000 --> 00:00:00,000 format)
		if timelineRegex.MatchString(trimmed) {
			continue
		}

		// Skip numeric lines (subtitle sequence numbers)
		if numericLinesRegex.MatchString(trimmed) {
			continue
		}

		// Replace special characters in lines (subtitle sequence numbers)
		if specialCharactersAtEndRegex.MatchString(trimmed) {
			trimmed = specialCharactersAtEndRegex.ReplaceAllString(trimmed, "")
		}

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Skip duplicate lines
		if seen[trimmed] {
			continue
		}

		// Write the line to output
		seen[trimmed] = true
		_, err := sb.WriteString(trimmed + " ")
		if err != nil {
			return "", err
		}
	}
	return strings.TrimSpace(sb.String()), nil
}
