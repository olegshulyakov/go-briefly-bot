package workers

type LoaderMessage struct {
	ClientAppID int8
	MessageID   int64
	UserID      int64
	Url         string
	Language    string
}

func (lm *LoaderMessage) Serialize() ([]byte, error) {
	// Serialize loader message to JSON
	return nil, nil
}

func DeserializeLoaderMessage(data []byte) (*LoaderMessage, error) {
	// Deserialize loader message from JSON
	return nil, nil
}
