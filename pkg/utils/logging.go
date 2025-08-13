package utils

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// LogLevel represents the severity of a log message.
// Using slog's levels: Debug(-4), Info(0), Warn(4), Error(8)
// We can alias them for clarity if needed, but using slog's constants is standard.
// type LogLevel slog.Level

// Logger wraps slog.Logger to provide a consistent interface.
type Logger struct {
	*slog.Logger
	appName string
}

// NewLogger creates a new structured logger instance using slog.
func NewLogger(config *Config) *Logger {
	// Determine log level from config
	var level slog.Level
	switch strings.ToLower(config.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		// Default to info level if invalid or not set
		level = slog.LevelInfo
		fmt.Printf("WARNING: Invalid or unspecified LOG_LEVEL '%s', defaulting to 'info'\n", config.LogLevel)
	}

	// Determine log format (handler)
	var handler slog.Handler
	switch strings.ToLower(config.LogFormat) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	case "text":
		fallthrough // Default to text
	default:
		// Text handler is good for development, JSON for production
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	// Create the base slog logger
	slogLogger := slog.New(handler)

	// Add common attributes, like app name
	appName := "Briefly" //TODO Get from config
	slogLogger = slogLogger.With(slog.String("app", appName))

	return &Logger{
		Logger:  slogLogger,
		appName: appName,
	}
}

// Fatal logs a fatal error message and exits the program.
// slog doesn't have a Fatal level, so we implement it using Error + os.Exit.
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.Logger.ErrorContext(context.Background(), msg, args...)
	os.Exit(1)
}
