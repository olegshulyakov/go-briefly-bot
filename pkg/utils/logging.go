package utils

import (
	"fmt"
	"log"
	"os"
)

// LogLevel represents the severity of a log message.
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of the LogLevel.
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger is a simple structured logger.
type Logger struct {
	level   LogLevel
	out     *log.Logger
	err     *log.Logger
	appName string
}

// NewLogger creates a new logger instance.
func NewLogger(config *Config) *Logger {
	logLevel := INFO
	switch config.LogLevel {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
	case "warn":
		logLevel = WARN
	case "error":
		logLevel = ERROR
	}

	// For simplicity, using standard log package.
	// A more advanced logger like Zap or Logrus could be used.
	outLogger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	errLogger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

	return &Logger{
		level:   logLevel,
		out:     outLogger,
		err:     errLogger,
		appName: "VideoSummaryBot", // Could be made configurable
	}
}

// logInternal handles the actual logging logic.
func (l *Logger) logInternal(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return // Don't log if below the configured level
	}

	var logger *log.Logger
	if level >= ERROR {
		logger = l.err
	} else {
		logger = l.out
	}

	levelStr := level.String()
	message := fmt.Sprintf(format, args...)
	// The Lshortfile flag will add the file and line number
	logger.Printf("[%s] [%s] %s", l.appName, levelStr, message)
}

// Debug logs a debug message.
func (l *Logger) Debug(format string, args ...interface{}) {
	l.logInternal(DEBUG, format, args...)
}

// Info logs an informational message.
func (l *Logger) Info(format string, args ...interface{}) {
	l.logInternal(INFO, format, args...)
}

// Warn logs a warning message.
func (l *Logger) Warn(format string, args ...interface{}) {
	l.logInternal(WARN, format, args...)
}

// Error logs an error message.
func (l *Logger) Error(format string, args ...interface{}) {
	l.logInternal(ERROR, format, args...)
}

// Fatal logs a fatal error message and exits the program.
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.logInternal(ERROR, format, args...)
	os.Exit(1)
}