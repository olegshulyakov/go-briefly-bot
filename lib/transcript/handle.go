// Package transcript provides functionality to retrieve video transcripts
// from YouTube videos. It extracts video metadata (title, uploader, thumbnail)
// along with the transcript text in the video's primary language.
package transcript

import (
	"fmt"

	"github.com/olegshulyakov/go-briefly-bot/lib/transcript/youtube"
)

// VideoTranscript represents the complete transcript information for a video,
// including metadata and the transcript text itself.
type VideoTranscript struct {
	ID         string `json:"id"`         // Unique YouTube video ID
	Language   string `json:"language"`   // Primary language code of the video (e.g., "en")
	Uploader   string `json:"uploader"`   // Name of the video uploader/channel
	Title      string `json:"title"`      // Title of the video
	Thumbnail  string `json:"thumbnail"`  // URL to the video thumbnail image
	Transcript string `json:"transcript"` // Full text transcript of the video
}

// Transcript retrieves and returns the complete transcript information
// for a YouTube video given its URL.
//
// Parameters:
//   - videoURL: The full YouTube video URL (e.g., "https://youtube.com/watch?v=...")
//
// Returns:
//   - *VideoTranscript containing all video metadata and transcript text
//   - error if any step of the process fails (video info or transcript retrieval)
//
// Example:
//
//	transcript, err := Transcript("https://youtube.com/watch?v=dQw4w9WgXcQ")
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(transcript.Title)       // Print video title
//	fmt.Println(transcript.Transcript)  // Print transcript text
func Transcript(videoURL string) (*VideoTranscript, error) {
	videoInfo, err := youtube.VideoInfo(videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get video info: %s", videoURL)
	}

	transcript, err := youtube.Transcript(videoURL, videoInfo.Language)
	if err != nil {
		return nil, fmt.Errorf("failed to get video transcript: %s", videoURL)
	}

	var videoTranscript = &VideoTranscript{
		ID:         videoInfo.ID,
		Uploader:   videoInfo.Uploader,
		Language:   videoInfo.Language,
		Title:      videoInfo.Title,
		Thumbnail:  videoInfo.Thumbnail,
		Transcript: transcript,
	}

	return videoTranscript, nil
}
