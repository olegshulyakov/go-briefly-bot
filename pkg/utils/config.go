package utils

type Config struct {
	DatabasePath  string
	APIServerPort string
	TelegramToken string
	OpenAIAPIKey  string
	RabbitMQURL   string
	BatchSizes    BatchSizes
	Timeouts      Timeouts
	RateLimits    RateLimits
}

type BatchSizes struct {
	LoaderProducer      int
	TransformerProducer int
	ResultHandler       int
}

type Timeouts struct {
	LoaderProducer      int
	TransformerProducer int
	ResultHandler       int
	RetryHandler        int
	ExpirationHandler   int
}

type RateLimits struct {
	UserRequestDelay int
	WarmupPeriod     int
}

func LoadConfig() (*Config, error) {
	// Load configuration from environment variables or config file
	return nil, nil
}
