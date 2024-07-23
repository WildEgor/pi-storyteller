package bot

import "image"

type StorySlide struct {
	Image image.Image
	Desc  string
}

type MessageRecipient struct {
	ID string
}

func (r *MessageRecipient) Recipient() string {
	return r.ID
}

type TelegramCommandData struct {
	ChatID  int64
	Payload string
}
