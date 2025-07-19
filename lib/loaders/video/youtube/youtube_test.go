package youtube_test

import (
	"testing"

	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/youtube"
)

func TestIsValidURL(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "valid full URL with http",
			input: "https://www.youtube.com/watch?v=abcdefghijk",
			valid: true,
		},
		{
			name:  "valid youtu.be",
			input: "youtu.be/12345678901",
			valid: true,
		},
		{
			name:  "valid URL with extra parameters",
			input: "http://youtube.com/watch?feature=shared&v=anotherID123",
			valid: true,
		},
		{
			name:  "valid URL without http",
			input: "www.youtube.com/watch?v=validID1234",
			valid: true,
		},
		{
			name:  "valid URL with text around",
			input: "text  https://youtu.be/validID1234  more text",
			valid: true,
		},
		{
			name:  "invalid: different domain",
			input: "https://vimeo.com/123",
			valid: false,
		},
		{
			name:  "invalid: short ID",
			input: "https://www.youtube.com/watch?v=short12",
			valid: false,
		},
		// {
		// 	name:  "invalid: long ID",
		// 	input: "youtu.be/123456789012",
		// 	valid: false,
		// },
		{
			name:  "invalid: embed URL",
			input: " https://youtube.com/embed/validID1234 ",
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := youtube.IsValidURL(tc.input); got != tc.valid {
				t.Errorf("IsValidURL(%q) = %v, want %v", tc.input, got, tc.valid)
			}
		})
	}
}

func TestGetID(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expected  string
		expectErr bool
	}{
		{
			name:      "empty input",
			input:     "",
			expectErr: true,
		},
		{
			name:      "valid URL with http",
			input:     "https://www.youtube.com/watch?v=abcdefghijk",
			expected:  "abcdefghijk",
			expectErr: false,
		},
		{
			name:      "valid youtu.be",
			input:     "youtu.be/12345678901",
			expected:  "12345678901",
			expectErr: false,
		},
		{
			name:      "valid URL with extra parameters",
			input:     "http://youtube.com/watch?feature=shared&v=validID1234",
			expected:  "validID1234",
			expectErr: false,
		},
		{
			name:      "valid URL with text around",
			input:     "text  https://youtu.be/validID1234  more text",
			expected:  "validID1234",
			expectErr: false,
		},
		{
			name:      "invalid: different domain",
			input:     "https://vimeo.com/123 ",
			expectErr: true,
		},
		{
			name:      "invalid: short ID",
			input:     "https://www.youtube.com/watch?v=short12",
			expectErr: true,
		},
		// {
		// 	name:      "invalid: long ID",
		// 	input:     "youtu.be/123456789012",
		// 	expectErr: true,
		// },
		{
			name:      "invalid: embed URL",
			input:     " https://youtube.com/embed/validID123 ",
			expectErr: true,
		},
		{
			name:      "invalid: missing v parameter",
			input:     "youtube.com/watch?feature=shared",
			expectErr: true,
		},
		{
			name:      "invalid: first v is invalid, second valid",
			input:     "youtube.com/watch?v=shortID&v=validID123",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := youtube.GetID(tc.input)
			if (err != nil) != tc.expectErr {
				t.Errorf("GetID(%q) got error: %v, expected error: %v", tc.input, err, tc.expectErr)
			}
			if !tc.expectErr && got != tc.expected {
				t.Errorf("GetID(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestExtractURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "multiple valid URLs",
			input:    "text https://youtu.be/validID1234  and youtube.com/watch?v=validID1234 more text",
			expected: []string{"https://youtu.be/validID1234", "youtube.com/watch?v=validID1234"},
		},
		{
			name:     "no URLs",
			input:    "no youtube URLs here",
			expected: []string{},
		},
		{
			name:     "URL with query parameters",
			input:    "https://www.youtube.com/watch?v=validID1234&list=playlist",
			expected: []string{"https://www.youtube.com/watch?v=validID1234"},
		},
		{
			name:     "mixed valid and invalid URLs",
			input:    "invalid.com/123 and youtu.be/validID1234 and  https://vimeo.com/456 ",
			expected: []string{"youtu.be/validID1234"},
		},
		{
			name:     "URLs with trailing text",
			input:    "https://youtu.be/validID1234?list=... and more",
			expected: []string{"https://youtu.be/validID1234"},
		},
		{
			name:     "multiple valid URLs with different formats",
			input:    "youtu.be/validID1234 https://www.youtube.com/watch?v=validID1234  https://youtu.be/validID1234 ",
			expected: []string{"youtu.be/validID1234", "https://www.youtube.com/watch?v=validID1234", "https://youtu.be/validID1234"},
		},
		{
			name:     "invalid URL with wrong ID length",
			input:    "youtu.be/shortID",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := youtube.ExtractURLs(tc.input)
			if len(got) != len(tc.expected) {
				t.Errorf("ExtractURLs(%q) got %d URLs, want %d", tc.input, len(got), len(tc.expected))
				return
			}
			for i := range got {
				if got[i] != tc.expected[i] {
					t.Errorf("ExtractURLs(%q) element %d: got %q, want %q", tc.input, i, got[i], tc.expected[i])
				}
			}
		})
	}
}
