package utils_test

import (
	"os"
	"strings"
	"testing"

	"github.com/olegshulyakov/go-briefly-bot/briefly/transcript/utils"
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
            expected: "Hello World ",
        },
        {
            name: "duplicates",
            input: `1
00:00:00,000 --> 00:00:00,001
Line
2
00:00:00,002 --> 00:00:00,003
Line`,
            expected: "Line ",
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
            expected: "   Text ",
        },
        {
            name: "empty_lines",
            input: `
Some text

Another line`,
            expected: "Some text Another line ",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            actual, err := utils.CleanSRT(tc.input)
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

// TestReadAndRemoveFile verifies the ReadAndRemoveFile function's behavior.
func TestReadAndRemoveFile(t *testing.T) {
    t.Run("valid file", func(t *testing.T) {
        content := "Test content"
        tmpFile, err := os.CreateTemp("", "testfile.*.txt")
        if err != nil {
            t.Fatal(err)
        }
        defer os.Remove(tmpFile.Name()) // Cleanup in case of test failure

        // Write content to the temporary file
        if _, err := tmpFile.Write([]byte(content)); err != nil {
            t.Fatal(err)
        }
        if err := tmpFile.Close(); err != nil {
            t.Fatal(err)
        }

        actualContent, err := utils.ReadAndRemoveFile(tmpFile.Name())
        if err != nil {
            t.Errorf("Unexpected error: %v", err)
        }
        if actualContent != content {
            t.Errorf("Expected %q, got %q", content, actualContent)
        }

        // Check if the file was deleted
        if _, err := os.Stat(tmpFile.Name()); !os.IsNotExist(err) {
            t.Errorf("File was not deleted")
        }
    })

    t.Run("empty filename", func(t *testing.T) {
        _, err := utils.ReadAndRemoveFile("")
        if err == nil {
            t.Error("Expected error for empty filename, got none")
        }
        if err.Error() != "file name is empty" {
            t.Errorf("Unexpected error message: %v", err)
        }
    })

    t.Run("non-existent file", func(t *testing.T) {
        _, err := utils.ReadAndRemoveFile("nonexistentfile.txt")
        if err == nil {
            t.Error("Expected error for non-existent file, got none")
        }
        if !strings.Contains(err.Error(), "no file found") {
            t.Errorf("Error message does not contain 'no file found': %v", err)
        }
    })
}