// Package video provides functionality to retrieve video transcripts
// from videos. It extracts video metadata (title, uploader, thumbnail)
// along with the transcript text in the video's primary language.
package video

// Info represents metadata about a video.
type Info struct {
	ID        string `json:"id"`        // The unique identifier of the video.
	Language  string `json:"language"`  // The video language.
	Uploader  string `json:"uploader"`  // The name of the video uploader.
	Title     string `json:"title"`     // The title of the video.
	Thumbnail string `json:"thumbnail"` // The URL of the video's thumbnail.
}

// Transcript represents the complete transcript information for a video,
// including metadata and the transcript text itself.
type Transcript struct {
	ID         string `json:"id"`         // Unique video ID
	Language   string `json:"language"`   // Primary language code of the video (e.g., "en")
	Uploader   string `json:"uploader"`   // Name of the video uploader/channel
	Title      string `json:"title"`      // Title of the video
	Thumbnail  string `json:"thumbnail"`  // URL to the video thumbnail image
	Transcript string `json:"transcript"` // Full text transcript of the video
}

// InfoLoader represents loader interface for a video metadata.
type InfoLoader interface {
	// VideoInfo returns metadata for the video.
	// Returns nil if metadata could not be retrieved.
	VideoInfo() *Info
}

// TranscriptLoader represents loader interface for a video transcript.
type TranscriptLoader interface {
	// Transcript returns the transcript text for the video.
	// Returns nil if the transcript could not be retrieved.
	Transcript() *Transcript
}

// DataLoader represents loader interface for a video.
type DataLoader interface {
	// Load initializes the video data. This method must be called
	// before accessing VideoInfo or Transcript. Returns an error if
	// data loading fails.
	Load() error

	// VideoInfo returns metadata for the video.
	// Returns nil if Load() was not called or failed.
	VideoInfo() *Info

	// Transcript returns the transcript text for the video.
	// Returns nil if Load() was not called or failed.
	Transcript() *Transcript
}
