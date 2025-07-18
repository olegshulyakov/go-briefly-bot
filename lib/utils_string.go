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
	if chunkSize <= 0 {
		return []string{} // Return the original string as a single chunk if chunkSize is not positive.
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
			chunks = append(chunks, currentChunk.String())
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
