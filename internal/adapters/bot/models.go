package bot

import (
	"image"
	"strconv"
)

type StorySlide struct {
	Image image.Image
	Desc  string
}

type MessageRecipient struct {
	ID        string
	ChatID    int64
	MessageID int64
}

func (r *MessageRecipient) Recipient() string {
	return r.ID
}

func (r *MessageRecipient) MessageSig() (messageID string, chatID int64) {
	if r.ChatID == 0 {
		v, _ := strconv.Atoi(r.ID)

		return strconv.Itoa(int(r.MessageID)), int64(v)
	}

	return strconv.Itoa(int(r.MessageID)), r.ChatID
}

type CommandData struct {
	Nickname string
	ChatID   int64
	Payload  string
}
