// Package ytdlp provides functionality for interacting with yt-dlp.
//
// The package uses yt-dlp (https://github.com/yt-dlp/yt-dlp) as a backend
// for fetching data from YouTube. Additional yt-dlp options can be configured
// through the YT_DLP_ADDITIONAL_OPTIONS environment variable.
package ytdlp

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

// Exec executes the yt-dlp command with the provided arguments and URL.
// It automatically retries the command up to 3 times on failure and includes
// any additional options specified in YT_DLP_ADDITIONAL_OPTIONS.
// Returns the command output or an error if all attempts fail.
func Exec(arguments []string, url string) ([]byte, error) {
	const maxAttempts = 3
	// ytDlpAdditionalOptions contains additional command-line options for yt-dlp
	// parsed from YT_DLP_ADDITIONAL_OPTIONS environment variable.
	var ytDlpAdditionalOptions = strings.Fields(os.Getenv("YT_DLP_ADDITIONAL_OPTIONS"))

	var (
		err    error
		output []byte
		args   = make([]string, 0, len(arguments)+len(ytDlpAdditionalOptions)+1)
	)

	if len(ytDlpAdditionalOptions) > 0 {
		args = append(args, ytDlpAdditionalOptions...)
	}
	args = append(args, arguments...)
	args = append(args, url)

	// Execute with retry
	slog.Debug("Executing yt-dlp", "args", args)
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if output, err = exec.Command("yt-dlp", args...).Output(); err == nil {
			break
		}
	}

	if err != nil {
		return output, fmt.Errorf("failed to extract transcript: %w\n%s", err, output)
	}

	return output, nil
}
