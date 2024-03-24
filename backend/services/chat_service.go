package services

// ChatMessage represents a message in a chat room
type ChatMessage struct {
	MessageType int
	Data        []byte
	Sender      string
}
