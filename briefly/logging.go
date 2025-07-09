package briefly

import (
	"log/slog"
	"os"
)

// Logger is a global logger instance used for logging throughout the application.
var Logger *slog.Logger = GetLogger()

// GetLogger initializes the global logger with a text formatter and sets the log level to Info.
//
// The logger is configured to output to standard output (stdout).
func GetLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	return logger
}
