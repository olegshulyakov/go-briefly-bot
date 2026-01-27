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
	if chunkSize <= 0 || text == "" || len(text) <= chunkSize {
		return []string{text} // Return the original string as a single chunk.
	}

	var (
		chunks       []string
		currentChunk strings.Builder
		runeCount    = 0
		runeIndex    = 0 // Track rune index within the text for slicing
	)

	for index := 0; index < len(text); {
		_, size := utf8.DecodeRuneInString(text[index:]) // Get the rune and its size in bytes

		if runeCount+1 > chunkSize {
			// Chunk is full, finalize and start a new one
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
			runeCount = 0
		}

		// Append the rune to the current chunk
		currentChunk.WriteString(text[index : index+size]) // Append the rune as bytes
		runeCount++
		index += size
		runeIndex++
	}

	// Append the last chunk if it's not empty
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
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
// func ToParagraphsAndChunks(text string, chunkSize int) []string {
// 	if chunkSize <= 0 || len(text) <= chunkSize {
// 		return []string{strings.TrimSpace(text)} // Return the original string as a single chunk if chunkSize is not positive.
// 	}

// 	// Split by paragraphs first
// 	paragraphs := strings.Split(text, "\n\n")

// 	var result []string

// 	for _, paragraph := range paragraphs {
// 		// Trim leading/trailing whitespace from each paragraph
// 		trimmedParagraph := strings.TrimSpace(paragraph)

// 		// If the paragraph is empty after trimming, skip it
// 		if trimmedParagraph == "" {
// 			continue
// 		}

// 		// If the paragraph length is within the chunk size, add it directly
// 		if utf8.RuneCountInString(trimmedParagraph) <= chunkSize {
// 			result = append(result, trimmedParagraph)
// 		} else {
// 			// If the paragraph is too long, split it using the existing ToChunks function
// 			chunks := ToChunks(trimmedParagraph, chunkSize)
// 			result = append(result, chunks...)
// 		}
// 	}

// 	return result
// }
