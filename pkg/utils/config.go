package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	// Database
	DatabaseBasePath        string
	DatabaseMaxOpenConns    int
	DatabaseMaxIdleConns    int
	DatabaseConnMaxLifetime time.Duration
	DatabaseConnMaxIdleTime time.Duration

	// API Server
	APIServerPort string
	APIKey        string

	// Telegram Bot
	TelegramToken          string
	TelegramRateLimitDelay time.Duration
	TelegramWarmupPeriod   time.Duration

	// OpenAI
	OpenAIAPIKey         string
	OpenAIBaseURL        string // For custom endpoints
	OpenAIModel          string
	OpenAISystemPrompt   string
	OpenAIChunkLimitSize int

	// RabbitMQ
	RabbitMQURL string

	// Worker Batches and Timeouts
	LoaderProducerBatchSize      int
	LoaderProducerTimeout        time.Duration
	TransformerProducerBatchSize int
	TransformerProducerTimeout   time.Duration
	ResultHandlerBatchSize       int
	ResultHandlerTimeout         time.Duration
	RetryHandlerRetryLimit       int
	RetryHandlerTimeout          time.Duration
	ExpirationHandlerTimeout     time.Duration

	// Logging
	LogLevel  string
	LogFormat string // json, text
}

// LoadConfig loads configuration from environment variables with defaults.
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// Database
	cfg.DatabaseBasePath = getEnv("DB_BASE_PATH", "./data")
	cfg.DatabaseMaxOpenConns = getEnvAsInt("DB_MAX_OPEN_CONNS", 10)
	cfg.DatabaseMaxIdleConns = getEnvAsInt("DB_MAX_IDLE_CONNS", 2)
	cfg.DatabaseConnMaxLifetime = getEnvAsDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute)
	// DB_CONN_MAX_IDLE_TIME is available in Go 1.15+
	// cfg.DatabaseConnMaxIdleTime = getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute)

	// API Server
	cfg.APIServerPort = getEnv("API_SERVER_PORT", "8080")
	cfg.APIKey = getEnv("API_KEY", "") // Should be required in production
	if cfg.APIKey == "" {
		// In a real system, this should probably be required or have a default secure value for dev
		fmt.Println("WARNING: API_KEY is not set")
	}

	// Telegram Bot
	cfg.TelegramToken = getEnv("TELEGRAM_TOKEN", "")
	cfg.TelegramRateLimitDelay = getEnvAsDuration("TELEGRAM_RATE_LIMIT_DELAY", 5*time.Second)
	cfg.TelegramWarmupPeriod = getEnvAsDuration("TELEGRAM_WARMUP_PERIOD", 30*time.Second)

	// OpenAI
	cfg.OpenAIAPIKey = getEnv("OPENAI_API_KEY", "")
	cfg.OpenAIBaseURL = getEnv("OPENAI_BASE_URL", "")         // Default to OpenAI's public API
	cfg.OpenAIModel = getEnv("OPENAI_MODEL", "gpt-3.5-turbo") // Or gpt-4
	cfg.OpenAISystemPrompt = getEnv("OPENAI_SYSTEM_PROMPT", "You are a helpful assistant that summarizes video transcripts.")
	cfg.OpenAIChunkLimitSize = getEnvAsInt("OPENAI_CHUNK_LIMIT_SIZE", 3000)

	// RabbitMQ
	cfg.RabbitMQURL = getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	// Worker Batches and Timeouts
	cfg.LoaderProducerBatchSize = getEnvAsInt("LOADER_PRODUCER_BATCH_SIZE", 10)
	cfg.LoaderProducerTimeout = getEnvAsDuration("LOADER_PRODUCER_TIMEOUT", 10*time.Second)
	cfg.TransformerProducerBatchSize = getEnvAsInt("TRANSFORMER_PRODUCER_BATCH_SIZE", 10)
	cfg.TransformerProducerTimeout = getEnvAsDuration("TRANSFORMER_PRODUCER_TIMEOUT", 10*time.Second)
	cfg.ResultHandlerBatchSize = getEnvAsInt("RESULT_HANDLER_BATCH_SIZE", 100)
	cfg.ResultHandlerTimeout = getEnvAsDuration("RESULT_HANDLER_TIMEOUT", 10*time.Second)
	cfg.RetryHandlerRetryLimit = getEnvAsInt("RETRY_HANDLER_RETRY_LIMIT", 5)
	cfg.RetryHandlerTimeout = getEnvAsDuration("RETRY_HANDLER_TIMEOUT", 30*time.Second)
	cfg.ExpirationHandlerTimeout = getEnvAsDuration("EXPIRATION_HANDLER_TIMEOUT", 4*time.Hour)

	// Logging
	cfg.LogLevel = getEnv("LOG_LEVEL", "info")
	cfg.LogFormat = getEnv("LOG_FORMAT", "text") // or "json"

	return cfg, nil
}

// Helper functions to get environment variables with defaults

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	fmt.Printf("WARNING: Invalid integer value for %s, using default %d\n", key, defaultValue)
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	// Handle common duration formats like "10s", "5m", "1h"
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	// Handle integer seconds as fallback
	if seconds, err := strconv.Atoi(valueStr); err == nil {
		return time.Duration(seconds) * time.Second
	}
	fmt.Printf("WARNING: Invalid duration value for %s, using default %v\n", key, defaultValue)
	return defaultValue
}

// getEnvAsBool parses a boolean environment variable.
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	// Convert to lowercase for comparison
	valueStr = strings.ToLower(valueStr)
	switch valueStr {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		fmt.Printf("WARNING: Invalid boolean value for %s, using default %v\n", key, defaultValue)
		return defaultValue
	}
}
