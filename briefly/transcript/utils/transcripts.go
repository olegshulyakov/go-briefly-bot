package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// CleanSRT cleans up a transcript string by removing empty lines,
// timeline lines, numeric lines (subtitle sequence numbers), and duplicate lines.
//
// It takes the transcript text as input and returns the cleaned transcript
// as a string and an error if any issue occurs during processing.
//
// The function performs the following cleaning operations:
//   - Skips empty lines: Lines consisting only of whitespace are removed.
//   - Skips timeline lines: Lines matching the SRT timeline format
//     "HH:MM:SS,milliseconds --> HH:MM:SS,milliseconds" are removed.
//   - Skips numeric lines: Lines that can be parsed as integers are removed,
//     assuming these are subtitle sequence numbers.
//   - Removes duplicate lines: Only the first occurrence of each unique line is kept.
//
// Example:
//
//	transcript := `1
//	00:00:00,000 --> 00:00:05,000
//	Hello world
//
//	2
//	00:00:05,000 --> 00:00:10,000
//	Hello world
//	Duplicate line
//
//	Duplicate line
//
//	3
//	00:00:10,000 --> 00:00:15,000
//	Another line
//	`
//	cleanedTranscript, err := CleanSRT(transcript)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//	fmt.Println(cleanedTranscript)
//	// Output:
//	// Hello world
//	// Duplicate line
//	// Another line
//
// Args:
//
//	text: The SRT transcript text as a string.
//
// Returns:
//
//	string: The cleaned transcript string.
//	error:  An error if there was an issue during processing. Returns nil in normal cases.
func CleanSRT(text string) (string, error) {
	var sb strings.Builder
	seen := make(map[string]bool)
	timelineRegex := regexp.MustCompile(`^\d{2}:\d{2}:\d{2},\d{3} --> \d{2}:\d{2}:\d{2},\d{3}$`)
	numericLinesRegex := regexp.MustCompile(`^\d+$`)
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Skip timeline lines (00:00:00,000 --> 00:00:00,000 format)
		if timelineRegex.MatchString(trimmed) {
			continue
		}

		// Skip numeric lines (subtitle sequence numbers)
		if numericLinesRegex.MatchString(trimmed) {
			continue
		}

		// Skip duplicate lines
		if seen[trimmed] {
			continue
		}

		// Write the line to output
		seen[trimmed] = true
		_, err := sb.WriteString(line + " ")
		if err != nil {
			return "", err
		}
	}
	return sb.String(), nil
}

func ReadAndRemoveFile(fileName string) (string, error) {
	if fileName == "" {
		return "", errors.New("file name is empty")
	}

	text, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("no file found: %w", err)
	}

	err = os.Remove(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to delete file: %w", err)
	}

	return string(text), nil
}
