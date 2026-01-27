package provider_test

import (
	"testing"

	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/provider"
)

func TestVkVideo_IsValidURL(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "valid full URL with http",
			input: "https://vkvideo.ru/video-12345678_123456789",
			valid: true,
		},
		{
			name:  "valid URL with extra parameters",
			input: "https://vkvideo.ru/video-12345678_123456789?feature=shared",
			valid: true,
		},
		{
			name:  "valid URL without http",
			input: "vkvideo.ru/video-12345678_123456789",
			valid: true,
		},
		{
			name:  "valid URL with text around",
			input: "text  https://vkvideo.ru/video-12345678_123456789  more text",
			valid: true,
		},
		{
			name:  "invalid: different domain",
			input: "https://video.com/123",
			valid: false,
		},
		{
			name:  "invalid ID",
			input: "https://vkvideo.ru/video-12345678",
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := provider.VkVideo.IsValidURL(tc.input); got != tc.valid {
				t.Errorf("IsValidURL(%q) = %v, want %v", tc.input, got, tc.valid)
			}
		})
	}
}

func TestVkVideo_GetID(t *testing.T) {
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
			input:     "https://vkvideo.ru/video-12345678_123456789",
			expected:  "video-12345678_123456789",
			expectErr: false,
		},
		{
			name:      "valid URL with extra parameters",
			input:     "https://vkvideo.ru/video-12345678_123456789?feature=shared",
			expected:  "video-12345678_123456789",
			expectErr: false,
		},
		{
			name:      "valid URL with text around",
			input:     "text  https://vkvideo.ru/video-12345678_123456789  more text",
			expected:  "video-12345678_123456789",
			expectErr: false,
		},
		{
			name:      "invalid: different domain",
			input:     "https://video.com/123 ",
			expectErr: true,
		},
		{
			name:      "invalid ID",
			input:     "https://vkvideo.ru/video-12345678",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := provider.VkVideo.GetID(tc.input)
			if (err != nil) != tc.expectErr {
				t.Errorf("GetID(%q) got error: %v, expected error: %v", tc.input, err, tc.expectErr)
			}
			if !tc.expectErr && got != tc.expected {
				t.Errorf("GetID(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestVkVideo_ExtractURLs(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "multiple valid URLs",
			input:    "text https://vkvideo.ru/video-12345678_123456789  and https://vkvideo.ru/video-1234_1234 more text",
			expected: []string{"https://vkvideo.ru/video-12345678_123456789", "https://vkvideo.ru/video-1234_1234"},
		},
		{
			name:     "no URLs",
			input:    "no URLs here",
			expected: []string{},
		},
		{
			name:     "URL with query parameters",
			input:    "https://vkvideo.ru/video-1234_1234&list=playlist",
			expected: []string{"https://vkvideo.ru/video-1234_1234"},
		},
		{
			name:     "mixed valid and invalid URLs",
			input:    "invalid.com/123 and vkvideo.ru/video-12345678_123456789 and  https://video.com/456 ",
			expected: []string{"vkvideo.ru/video-12345678_123456789"},
		},
		{
			name:     "URLs with trailing text",
			input:    "https://vkvideo.ru/video-12345678_123456789?list=... and more",
			expected: []string{"https://vkvideo.ru/video-12345678_123456789"},
		},
		{
			name:     "multiple valid URLs with different formats",
			input:    "vkvideo.ru/video-1234_1234  https://vkvideo.ru/video-12345678_123456789 ",
			expected: []string{"vkvideo.ru/video-1234_1234", "https://vkvideo.ru/video-12345678_123456789"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := provider.VkVideo.ExtractURLs(tc.input)
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
