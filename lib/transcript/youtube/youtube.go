package youtube

import (
	"errors"
	"regexp"
)

// Constant string representing the regular expression pattern used to match YouTube video URLs.
//   - This pattern is designed to handle various YouTube URL formats, including:
//   - - Full URLs with "https://" or "http://"
//   - - URLs with "www."
//   - - Shortened URLs like "youtu.be/"
//   - - URLs with the standard "youtube.com/watch?v="
const youtubeURLPattern = `(?:https?:\/\/)?(?:www\.)?(?:youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]{11})`

// IsValidYouTubeURL checks if the provided text contains a valid YouTube URL.
//
// It compiles a regular expression defined by YoutubeUrlPattern and uses it to
// determine whether the input text contains a YouTube URL.
//
// Parameters:
//   - text: The string to check for a YouTube URL.
//
// Returns:
//   - true if the text contains a valid YouTube URL, false otherwise.
func IsValidYouTubeURL(text string) bool {
	re := getYouTubeRegexp()
	return re.MatchString(text)
}

func GetYouTubeID(text string) (string, error) {
	if !IsValidYouTubeURL(text) {
		return "", errors.New("no valid URL found")
	}

	re := getYouTubeRegexp()
	matches := re.FindStringSubmatch(text)
	if len(matches) == 0 {
		return "", errors.New("YouTube ID not found")
	}
	return matches[1], nil
}

// ExtractAllYouTubeURLs extracts all YouTube URLs from the given text.
// It uses a regular expression to find all matching URLs.
//
// Parameters:
//   - text: The string to extract YouTube URLs from.
//
// Returns:
//   - A slice of strings containing all the YouTube URLs found in the text.
//   - An error if there was an error compiling the regex.
//
// Example:
//
//	urls, err := ExtractAllYouTubeURLs("Check out this video: https://www.youtube.com/watch?v=dQw4w9WgXcQ and another one at https://youtu.be/abcdefg123")
//	if err != nil {
//	    log.Errorf("Error extracting URLs: %v", err)
//	}
//	fmt.Println("URLs:", urls)
func ExtractAllYouTubeURLs(text string) ([]string, error) {
	if !IsValidYouTubeURL(text) {
		return nil, errors.New("no valid URL found")
	}

	re := getYouTubeRegexp()
	return re.FindAllString(text, -1), nil
}

func getYouTubeRegexp() *regexp.Regexp {
	return regexp.MustCompile(youtubeURLPattern)
}
