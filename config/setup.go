package config

// SetupConfig initializes full configuration
//
// Example:
//
// SetupConfig()
func SetupConfig() (*Config, error) {
	// Set up logger
	setupLogger()

	// Set up localizer
	setupLocalizer()

	// Load configuration
	return loadConfig()
}
