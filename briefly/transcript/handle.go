package transcript

import (
	"fmt"

	"github.com/olegshulyakov/go-briefly-bot/briefly/transcript/youtube"
)

// VideoTranscript represents metadata about a YouTube video transcript.
type VideoTranscript struct {
	ID         string `json:"id"`         // The unique identifier of the video.
	Uploader   string `json:"uploader"`   // The name of the video uploader.
	Language   string `json:"language"`   // The video language.
	Title      string `json:"title"`      // The title of the video.
	Thumbnail  string `json:"thumbnail"`  // The URL of the video's thumbnail.
	Transcript string `json:"transcript"` // The URL of the video's transcript.
}

func GetYoutubeVideoTranscript(videoURL string) (*VideoTranscript, error) {
	videoInfo, err := youtube.GetYoutubeVideoInfo(videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get video info: %s", videoURL)
	}

	transcript, err := youtube.GetYoutubeTranscript(videoURL, videoInfo.Language)
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
