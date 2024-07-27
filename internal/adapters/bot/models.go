// Package bot represent abstraction for bot
package bot

import (
	"context"
	"image"
	"strconv"
)

// Bot ...
type Bot interface {
	SendStory(ctx context.Context, to *MessageRecipient, slides []StorySlide) error
	SendMsg(ctx context.Context, to *MessageRecipient, msg string) (mid int, err error)
	EditMsg(ctx context.Context, to *MessageRecipient, msg string) (mid int, err error)
	DeleteMsg(ctx context.Context, to *MessageRecipient) error
	Start()
	Stop()
}

// Registry holds handlers
type Registry interface {
	HandleCommand(ctx context.Context, command string, fn commandHandler)
	RegisterBtnCallback(ctx context.Context, name string, fn btnCallbackHandler)
}

// CommandData ...
type CommandData struct {
	Nickname  string
	MessageID int
	ChatID    int64
	Lang      string
	Payload   string
}

// BtnCallbackData ...
type BtnCallbackData struct {
	Nickname  string
	MessageID int
	ChatID    int64
	Payload   string
}

type commandHandler func(data *CommandData) error
type btnCallbackHandler func(data *BtnCallbackData) error

// StorySlide ...
type StorySlide struct {
	ID    int
	Image image.Image
	Style string
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
		//nolint
		v, _ := strconv.Atoi(r.ID)
		return r.MessageID, int64(v)
	}
	return r.MessageID, r.ChatID
}
