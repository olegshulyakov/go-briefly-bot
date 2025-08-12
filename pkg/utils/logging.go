package utils

import "log"

func SetupLogger() {
	// Setup structured logging
}

func Info(message string, args ...interface{}) {
	// Log info message
	log.Printf(message, args...)
}

func Error(message string, args ...interface{}) {
	// Log error message
	log.Printf("ERROR: "+message, args...)
}

func Debug(message string, args ...interface{}) {
	// Log debug message
	log.Printf("DEBUG: "+message, args...)
}
