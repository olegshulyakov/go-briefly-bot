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
			expected:  []string{},
		},
		{
			name:      "negative chunk size",
			text:      "test",
			chunkSize: -1,
			expected:  []string{},
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
			expected:  []string{},
		},
		{
			name:      "empty string with zero chunk size",
			text:      "",
			chunkSize: 0,
			expected:  []string{},
		},
		{
			name:      "single rune chunk 0",
			text:      "x",
			chunkSize: 0,
			expected:  []string{},
		},
		{
			name:      "empty string chunk 1",
			text:      "",
			chunkSize: 1,
			expected:  []string{},
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
			name:      "multibyte characters split correctly",
			text:      "你好世界",
			chunkSize: 2,
			expected:  []string{"你好", "世界"},
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
