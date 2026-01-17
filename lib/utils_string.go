package lib

import (
	"strings"
	"unicode/utf8"
)

// ToChunks splits a string into chunks of specified rune count.
//
// The function takes a string and a chunkSize parameter, and returns a slice of strings
// where each string (except possibly the last one) contains exactly chunkSize runes.
// If chunkSize is less than or equal to 0, the empty chunk.
//
// This function properly handles UTF-8 encoded strings, counting runes rather than bytes,
// ensuring that multi-byte characters are not split across chunks.
//
// Example:
//
//	chunks := ToChunks("Hello, 世界", 3)
//	// chunks will be []string{"Hel", "lo,", " 世", "界"}
//
// Parameters:
//   - text: The input string to be split into chunks
//   - chunkSize: The maximum number of runes per chunk
//
// Returns:
//   - A slice of strings containing the chunks.
func ToChunks(text string, chunkSize int) []string {
	// Handle edge cases
	if chunkSize <= 0 {
		return []string{text}
	}
	if text == "" {
		return []string{""}
	}

	var chunks []string
	var current string

	for _, word := range strings.Fields(text) {
		// Handle long words
		if utf8.RuneCountInString(word) > chunkSize {
			if current != "" {
				chunks = append(chunks, current)
				current = ""
			}

			// Split long word
			runes := []rune(word)
			for i := 0; i < len(runes); i += chunkSize {
				end := min(i+chunkSize, len(runes))
				chunks = append(chunks, string(runes[i:end]))
			}
			continue
		}

		// Try to add word to current chunk
		withSpace := ""
		if current != "" {
			withSpace = current + " "
		}
		if utf8.RuneCountInString(withSpace+word) <= chunkSize {
			current = withSpace + word
		} else {
			if current != "" {
				chunks = append(chunks, current)
			}
			current = word
		}
	}

	if current != "" {
		chunks = append(chunks, current)
	}

	return chunks
}

// ToParagraphsAndChunks splits a string by paragraphs first (\n\n), then by chunk size if needed.
//
// The function takes a string and a chunkSize parameter, and returns a slice of strings
// where each string contains paragraphs or chunks of text. Paragraphs are defined as text
// separated by double newlines (\n\n). If a paragraph exceeds the chunkSize, it will be
// further divided into chunks of the specified size.
//
// Example:
//
//	result := ToChunksAndParagraphs("Hello Ben,\n\nHow are you?", 20)
//	// result will be []string{"Hello Ben,", "How are you?"}
//
// Parameters:
//   - text: The input string to be split by paragraphs and chunks
//   - chunkSize: The maximum number of runes per chunk (when needed)
//
// Returns:
//   - A slice of strings containing the paragraphs or chunks.
func ToParagraphsAndChunks(text string, chunkSize int) []string {
	if chunkSize <= 0 || len(text) <= chunkSize {
		return []string{strings.TrimSpace(text)} // Return the original string as a single chunk if chunkSize is not positive.
	}

	// Split by paragraphs first
	paragraphs := strings.Split(text, "\n\n")

	var result []string

	for _, paragraph := range paragraphs {
		// Trim leading/trailing whitespace from each paragraph
		trimmedParagraph := strings.TrimSpace(paragraph)

		// If the paragraph is empty after trimming, skip it
		if trimmedParagraph == "" {
			continue
		}

		// If the paragraph length is within the chunk size, add it directly
		if utf8.RuneCountInString(trimmedParagraph) <= chunkSize {
			result = append(result, trimmedParagraph)
		} else {
			// If the paragraph is too long, split it using the existing ToChunks function
			chunks := ToChunks(trimmedParagraph, chunkSize)
			result = append(result, chunks...)
		}
	}

	return result
}
