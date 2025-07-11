package briefly_test

import (
	"os"
	"testing"

	"github.com/olegshulyakov/go-briefly-bot/briefly"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resetConfig clears the global Configuration and environment variables
func resetConfig(envVars []string) {
	for _, v := range envVars {
		os.Unsetenv(v)
	}
}

func TestLoadConfiguration_SuccessFromEnvironment(t *testing.T) {
	// Define test environment variables
	envVars := []string{
		"TELEGRAM_BOT_TOKEN",
		"YT_DLP_ADDITIONAL_OPTIONS",
		"OPENAI_BASE_URL",
		"OPENAI_API_KEY",
		"OPENAI_MODEL",
	}
	defer resetConfig(envVars)

	// Set environment variables
	t.Setenv("TELEGRAM_BOT_TOKEN", "tele_token_env")
	t.Setenv("YT_DLP_ADDITIONAL_OPTIONS", "--option1  val1   --option2  val2")
	t.Setenv("OPENAI_BASE_URL", "openai_url_env")
	t.Setenv("OPENAI_API_KEY", "openai_key_env")
	t.Setenv("OPENAI_MODEL", "openai_model_env")

	// Execute and validate
	cfg, err := briefly.LoadConfiguration()
	require.NoError(t, err, "Valid config should not return error")
	assert.Equal(t, &briefly.Config{
		TelegramToken:          "tele_token_env",
		YtDlpAdditionalOptions: []string{"--option1", "val1", "--option2", "val2"},
		OpenAiBaseURL:          "openai_url_env",
		OpenAiAPIKey:           "openai_key_env",
		OpenAiModel:            "openai_model_env",
	}, cfg, "Configuration should match environment values")
}

func TestLoadConfiguration_ValidationFailures(t *testing.T) {
	tests := []struct {
		name        string
		envSetup    func()
		expectedErr string
	}{
		{
			name: "Missing Telegram Token",
			envSetup: func() {
				t.Setenv("OPENAI_BASE_URL", "url")
				t.Setenv("OPENAI_API_KEY", "key")
				t.Setenv("OPENAI_MODEL", "model")
			},
			expectedErr: "TELEGRAM_BOT_TOKEN not set",
		},
		{
			name: "Missing OpenAI URL",
			envSetup: func() {
				t.Setenv("TELEGRAM_BOT_TOKEN", "token")
				t.Setenv("OPENAI_API_KEY", "key")
				t.Setenv("OPENAI_MODEL", "model")
			},
			expectedErr: "OPENAI_BASE_URL not set",
		},
		{
			name: "Missing OpenAI Key",
			envSetup: func() {
				t.Setenv("TELEGRAM_BOT_TOKEN", "token")
				t.Setenv("OPENAI_BASE_URL", "url")
				t.Setenv("OPENAI_MODEL", "model")
			},
			expectedErr: "OPENAI_API_KEY not set",
		},
		{
			name: "Missing OpenAI Model",
			envSetup: func() {
				t.Setenv("TELEGRAM_BOT_TOKEN", "token")
				t.Setenv("OPENAI_BASE_URL", "url")
				t.Setenv("OPENAI_API_KEY", "key")
			},
			expectedErr: "OPENAI_MODEL not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset configuration and environment
			resetConfig([]string{
				"TELEGRAM_BOT_TOKEN",
				"OPENAI_BASE_URL",
				"OPENAI_API_KEY",
				"OPENAI_MODEL",
			})

			// Setup specific environment for test case
			tt.envSetup()

			// Execute and validate error
			_, err := briefly.LoadConfiguration()
			require.Error(t, err, "Expected validation error")
			assert.EqualError(t, err, tt.expectedErr, "Error message should match")
		})
	}
}

func TestLoadConfiguration_OptionalYtDlpOptions(t *testing.T) {
	// Define test environment variables
	envVars := []string{
		"TELEGRAM_BOT_TOKEN",
		"YT_DLP_ADDITIONAL_OPTIONS",
		"OPENAI_BASE_URL",
		"OPENAI_API_KEY",
		"OPENAI_MODEL",
	}
	defer resetConfig(envVars)

	// Set required variables only
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("OPENAI_BASE_URL", "url")
	t.Setenv("OPENAI_API_KEY", "key")
	t.Setenv("OPENAI_MODEL", "model")

	// Execute and validate
	cfg, err := briefly.LoadConfiguration()
	require.NoError(t, err, "Config with optional field missing should succeed")
	assert.Empty(t, cfg.YtDlpAdditionalOptions, "Optional field should be empty")
}
