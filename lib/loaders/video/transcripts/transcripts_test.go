package transcripts_test

import (
	"testing"

	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/transcripts"
)

// TestCleanSRT verifies that the CleanSRT function correctly processes SRT text.
func TestCleanSRT(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "basic case",
			input: `1
00:00:00,000 --> 00:00:00,001
Hello
2
00:00:00,002 --> 00:00:00,003
World`,
			expected: "Hello World",
		},
		{
			name: "duplicates",
			input: `1
00:00:00,000 --> 00:00:00,001
Line
2
00:00:00,002 --> 00:00:00,003
Line`,
			expected: "Line",
		},
		{
			name: "all_skipped",
			input: `1
00:00:00,000 --> 00:00:00,001
2`,
			expected: "",
		},
		{
			name: "leading_and_trailing_spaces",
			input: `   Text
 Text  `,
			expected: "Text",
		},
		{
			name: "special_symbols",
			input: `
Some text\h\h

Another line\h`,
			expected: "Some text Another line",
		},
		{
			name: "empty_lines",
			input: `
Some text

Another line`,
			expected: "Some text Another line",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := transcripts.CleanSRT(tc.input)
			if err != nil {
				t.Errorf("Unexpected error in test %s: %v", tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %s: Expected %q, got %q", tc.name, tc.expected, actual)
			}
		})
	}
}
