package db

type DictClientApp struct {
	ID  int8
	App string
}

type ProcessingStatus struct {
	ID          int8
	DisplayName string
}

type MessageHistory struct {
	ClientAppID    int8
	MessageID      int64
	UserID         int64
	UserName       string
	UserLanguage   string
	MessageContent string
	CreatedAt      string
}

type ProcessingQueue struct {
	ClientAppID  int8
	MessageID    int64
	UserID       int64
	Url          string
	Language     string
	StatusID     int8
	CreatedAt    string
	ProcessedAt  *string
	RetryCount   *int8
	ErrorMessage *string
}

type Source struct {
	Url       string
	Language  string
	Title     string
	Text      string
	CreatedAt string
}

type Summary struct {
	Url       string
	Language  string
	Summary   string
	CreatedAt string
}
