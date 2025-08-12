package workers

type TransformerMessage struct {
	ClientAppID int8
	MessageID   int64
	UserID      int64
	Url         string
	Language    string
	Text        string
}

func (tm *TransformerMessage) Serialize() ([]byte, error) {
	// Serialize transformer message to JSON
	return nil, nil
}

func DeserializeTransformerMessage(data []byte) (*TransformerMessage, error) {
	// Deserialize transformer message from JSON
	return nil, nil
}
