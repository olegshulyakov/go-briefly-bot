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
			chunks = appendWithTrimSpace(chunks, currentChunk.String())
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
		chunks = appendWithTrimSpace(chunks, currentChunk.String())
		currentChunk.Reset()
	}

	return chunks
}

// ToLexicalChunks splits a text string into semantically meaningful chunks of approximately
// the specified size, prioritizing natural language boundaries over arbitrary cuts.
//
// The function respects this priority order for break points:
//   - Paragraph boundaries (newlines) are preferred
//   - Sentence boundaries (., !, ?) are next
//   - Word boundaries (spaces) are last resort
//
// If no natural break point exists within chunkSize characters, the chunk
// will be cut at exactly chunkSize characters, even mid-word.
//
// Example:
//
//	text := "First sentence. Second!\n\nNew paragraph here."
//	chunks := ToLexicalChunks(text, 25)
//
// Parameters:
//   - text: The input string to be split into chunks
//   - chunkSize: The maximum number of runes per chunk
//
// Returns:
//
//	A slice of strings containing the lexical chunks.
func ToLexicalChunks(text string, chunkSize int) []string {
	if chunkSize <= 0 || len(text) <= chunkSize {
		return []string{strings.TrimSpace(text)} // Return the original string as a single chunk if chunkSize is not positive.
	}

	var chunks []string
	runes := []rune(text)

	for left := 0; left < len(runes); {
		right := min(left+chunkSize, len(runes))

		// If we're at the end of the text
		if right == len(runes) {
			currentChunk := string(runes[left:right])
			chunks = appendWithTrimSpace(chunks, currentChunk)
			break
		}

		right = findNaturalBreakpoint(runes, left, right)

		// Ensure we don't go past the text length
		right = min(right, len(runes))

		// Append to chunk array
		currentChunk := string(runes[left:right])
		chunks = appendWithTrimSpace(chunks, currentChunk)
		left = right

		// Skip whitespace at the beginning of next chunk
		for left < len(runes) && (runes[left] == ' ' || runes[left] == '\n') {
			left++
		}
	}

	return chunks
}

func appendWithTrimSpace(arr []string, text string) []string {
	cleaned := strings.TrimSpace(text)
	if len(cleaned) == 0 {
		return arr
	}
	return append(arr, cleaned)
}

// findNaturalBreakpoint finds natural right break point
// 1. Find index of end of last paragraph (newline), if we are at the end of paragraph - do nothing
// 2. Find index of end of last sentence
// 3. Find index of end of last word (last whitespace before boundary), if we are at the end of word - do nothing
// If no word break found, use the original right boundary.
func findNaturalBreakpoint(runes []rune, left int, right int) int {
	if runes[right] == '\n' {
		return right
	}

	textPiece := string(runes[left:right])

	lastParagraphIdx := strings.LastIndex(textPiece, "\n")
	if lastParagraphIdx > 0 {
		return left + lastParagraphIdx + 1 // Include the newline
	}

	lastSentenceIdx := max(strings.LastIndex(textPiece, "."), strings.LastIndex(textPiece, "!"), strings.LastIndex(textPiece, "?"))
	if lastSentenceIdx > 0 {
		return left + lastSentenceIdx + 1 // Include the punctuation
	}

	if runes[right] != ' ' {
		lastWordIdx := strings.LastIndex(textPiece, " ")
		if lastWordIdx > 0 {
			return left + lastWordIdx + 1
		}
	}

	return right
}
