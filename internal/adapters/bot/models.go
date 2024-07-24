package bot

import (
	"image"
	"strconv"
)

// StorySlide ...
type StorySlide struct {
	ID    int
	Image image.Image
	Desc  string
}

// MessageRecipient ...
type MessageRecipient struct {
	ID        string
	ChatID    int64
	MessageID string
}

// Recipient ...
func (r *MessageRecipient) Recipient() string {
	return r.ID
}

// MessageSig ...
func (r *MessageRecipient) MessageSig() (messageID string, chatID int64) {
	if r.ChatID == 0 {
		v, _ := strconv.Atoi(r.ID)
		return r.MessageID, int64(v)
	}
	return r.MessageID, r.ChatID
}

// CommandData ...
type CommandData struct {
	Nickname  string
	MessageID int
	ChatID    int64
	Payload   string
}
