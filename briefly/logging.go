package briefly

import (
	"log/slog"
	"os"
	"sync/atomic"
)

// logger is a global logger instance used for logging throughout the application.
var defaultLogger atomic.Pointer[slog.Logger]

func init() {
	defaultLogger.Store(setupLogger())
}

// Default returns the default [Logger].
func getLogger() *slog.Logger { return defaultLogger.Load() }

// setupLogger initializes the global logger with a text formatter and sets the log level to Info.
//
// The logger is configured to output to standard output (stdout).
func setupLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	return logger
}

func Debug(msg string, args ...any) {
	getLogger().Debug(msg, args...)
}

func Info(msg string, args ...any) {
	getLogger().Info(msg, args...)
}

func Warn(msg string, args ...any) {
	getLogger().Warn(msg, args...)
}

func Error(msg string, args ...any) {
	getLogger().Error(msg, args...)
}
