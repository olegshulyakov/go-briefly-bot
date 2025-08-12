package db

// These structs represent the database tables.
// They are used for serialization/deserialization in Go code.

type DictClientApp struct {
	ID  int8   `json:"id"`
	App string `json:"app"`
}

type ProcessingStatus struct {
	ID          int8   `json:"id"`
	DisplayName string `json:"display_name"`
}

type MessageHistory struct {
	ClientAppID    int8   `json:"client_app_id"`
	MessageID      int64  `json:"message_id"`
	UserID         int64  `json:"user_id"`
	UserName       string `json:"user_name"`
	UserLanguage   string `json:"user_language"`
	MessageContent string `json:"message_content"`
	CreatedAt      string `json:"created_at"` // Store as string for now, can be time.Time
}

type ProcessingQueue struct {
	ClientAppID  int8    `json:"client_app_id"`
	MessageID    int64   `json:"message_id"`
	UserID       int64   `json:"user_id"`
	Url          string  `json:"url"`
	Language     string  `json:"language"`
	StatusID     int8    `json:"status_id"`
	CreatedAt    string  `json:"created_at"`
	ProcessedAt  *string `json:"processed_at,omitempty"`  // Nullable
	RetryCount   *int8   `json:"retry_count,omitempty"`   // Nullable
	ErrorMessage *string `json:"error_message,omitempty"` // Nullable
}

type Source struct {
	Url       string `json:"url"`
	Language  string `json:"language"`
	Title     string `json:"title"`
	Text      string `json:"text"` // TEXT column
	CreatedAt string `json:"created_at"`
}

type Summary struct {
	Url       string `json:"url"`
	Language  string `json:"language"`
	Summary   string `json:"summary"` // TEXT column
	CreatedAt string `json:"created_at"`
}
