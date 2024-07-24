// Package bot represent abstraction for bot
package bot

import "context"

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
	HandleCommand(ctx context.Context, command string, fn func(data *CommandData) error)
}
