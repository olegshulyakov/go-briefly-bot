package lib_test

import (
	"testing"

	"github.com/olegshulyakov/go-briefly-bot/lib"
)

func TestToChunks(t *testing.T) {
	testCases := []struct {
		name      string
		text      string
		chunkSize int
		expected  []string
	}{
		{
			name:      "zero chunk size",
			text:      "hello",
			chunkSize: 0,
			expected:  []string{"hello"},
		},
		{
			name:      "negative chunk size",
			text:      "test",
			chunkSize: -1,
			expected:  []string{"test"},
		},
		{
			name:      "chunk larger than text",
			text:      "abc",
			chunkSize: 5,
			expected:  []string{"abc"},
		},
		{
			name:      "exact split",
			text:      "abcdef",
			chunkSize: 2,
			expected:  []string{"ab", "cd", "ef"},
		},
		{
			name:      "empty string with positive chunk size",
			text:      "",
			chunkSize: 5,
			expected:  []string{""},
		},
		{
			name:      "empty string with zero chunk size",
			text:      "",
			chunkSize: 0,
			expected:  []string{""},
		},
		{
			name:      "single rune chunk 0",
			text:      "x",
			chunkSize: 0,
			expected:  []string{"x"},
		},
		{
			name:      "empty string chunk 1",
			text:      "",
			chunkSize: 1,
			expected:  []string{""},
		},
		{
			name:      "single rune chunk",
			text:      "a",
			chunkSize: 1,
			expected:  []string{"a"},
		},
		{
			name:      "chunk size 1 with multiple runes",
			text:      "abcd",
			chunkSize: 1,
			expected:  []string{"a", "b", "c", "d"},
		},
		{
			name:      "mixed single and multi-byte",
			text:      "a😊b",
			chunkSize: 2,
			expected:  []string{"a😊", "b"},
		},
		{
			name:      "remainder chunk",
			text:      "abcdefg",
			chunkSize: 3,
			expected:  []string{"abc", "def", "g"},
		},
		{
			name:      "Paragraph exceeding chunk size",
			text:      "This is a very long paragraph that exceeds the chunk size and should be split into multiple chunks.",
			chunkSize: 10,
			expected: []string{
				"This is a",
				"very long",
				"paragraph",
				"that",
				"exceeds",
				"the chunk",
				"size and",
				"should be",
				"split into",
				"multiple",
				"chunks.",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := lib.ToChunks(tc.text, tc.chunkSize)
			if len(got) != len(tc.expected) {
				t.Errorf("For %v, expected %v chunks, got %v", tc.name, len(tc.expected), len(got))
			}
			for i := range got {
				if got[i] != tc.expected[i] {
					t.Errorf("For %v, chunk %d: expected %q, got %q", tc.name, i+1, tc.expected[i], got[i])
				}
			}
		})
	}
}

func TestToParagraphsAndChunks(t *testing.T) {
	testCases := []struct {
		name      string
		text      string
		chunkSize int
		expected  []string
	}{
		{
			name:      "Basic paragraph split",
			text:      "Hello Ben,\n\nHow are you?",
			chunkSize: 20,
			expected:  []string{"Hello Ben,", "How are you?"},
		},
		{
			name:      "Single paragraph within chunk size",
			text:      "Hello world!",
			chunkSize: 20,
			expected:  []string{"Hello world!"},
		},
		{
			name:      "Single paragraph exceeding chunk size",
			text:      "This is a very long paragraph that exceeds the chunk size and should be split into multiple chunks.",
			chunkSize: 10,
			expected: []string{
				"This is a",
				"very long",
				"paragraph",
				"that",
				"exceeds",
				"the chunk",
				"size and",
				"should be",
				"split into",
				"multiple",
				"chunks.",
			},
		},
		{
			name:      "Multiple paragraphs",
			text:      "Short.\n\nThis is a longer paragraph that will need to be split into smaller chunks.\n\nAnother short one.",
			chunkSize: 80,
			expected:  []string{"Short.", "This is a longer paragraph that will need to be split into smaller chunks.", "Another short one."},
		},
		{
			name:      "Multiple paragraphs with mixed sizes",
			text:      "Short.\n\nThis is a longer paragraph that will need to be split into smaller chunks.\n\nAnother short one.",
			chunkSize: 15,
			expected: []string{
				"Short.",
				"This is a",
				"longer",
				"paragraph that",
				"will need to be",
				"split into",
				"smaller chunks.",
				"Another short",
				"one.",
			},
		},
		{
			name:      "Empty text",
			text:      "",
			chunkSize: 10,
			expected:  []string{""},
		},
		{
			name:      "Only newlines",
			text:      "\n\n\n\n",
			chunkSize: 10,
			expected:  []string{""},
		},
		{
			name:      "Text with extra whitespace around paragraphs",
			text:      "  Hello Ben,  \n\n  How are you?  ",
			chunkSize: 20,
			expected:  []string{"Hello Ben,", "How are you?"},
		},
		{
			name:      "Zero chunk size",
			text:      "Hello world",
			chunkSize: 0,
			expected:  []string{"Hello world"},
		},
		{
			name:      "Negative chunk size",
			text:      "Hello world",
			chunkSize: -5,
			expected:  []string{"Hello world"},
		},
		{
			name:      "UTF-8 characters",
			text:      "Hello Ben,\n\nHow are you?",
			chunkSize: 10,
			expected:  []string{"Hello Ben,", "How are", "you?"},
		},
		{
			name:      "UTF-8 characters",
			text:      "Hello Ben,\n\nHow are you?",
			chunkSize: 50,
			expected:  []string{"Hello Ben,\n\nHow are you?"},
		},
		{
			name:      "Long UTF-8 paragraph",
			text:      "Hello Ben, this is a paragraph with UTF-8 characters that needs to be split properly.",
			chunkSize: 10,
			expected:  []string{"Hello Ben,", "this is a", "paragraph", "with UTF-8", "characters", "that needs", "to be", "split", "properly."},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := lib.ToParagraphsAndChunks(tc.text, tc.chunkSize)
			if len(got) != len(tc.expected) {
				t.Errorf("For %v, expected %v chunks, got %v", tc.name, len(tc.expected), len(got))
			}
			for i := range got {
				if got[i] != tc.expected[i] {
					t.Errorf("For %v, chunk %d: expected %q, got %q", tc.name, i+1, tc.expected[i], got[i])
				}
			}
		})
	}
}
