package lib

import (
	"strings"
	"unicode/utf8"
)

// SplitStringIntoChunks splits a string into chunks of a maximum length,
// ensuring that each chunk is valid UTF-8 and doesn't break runes (characters).
func SplitStringIntoChunks(text string, chunkSize int) []string {
	if chunkSize <= 0 {
		return []string{text} // Return the original string as a single chunk if chunkSize is not positive.
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
