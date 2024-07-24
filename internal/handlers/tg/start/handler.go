package tg_start_handler

import (
	"context"
	"fmt"

	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
)

// StartHandler ...
type StartHandler struct {
	bot bot.Bot
}

// NewStartHandler ...
func NewStartHandler(bot bot.Bot) *StartHandler {
	return &StartHandler{
		bot,
	}
}

// Handle ...
func (h *StartHandler) Handle(ctx context.Context, payload *StartPayload) error {
	_, err := h.bot.SendMsg(ctx, &bot.MessageRecipient{
		ID: payload.ChatID,
	}, fmt.Sprintf("Hi, %s", payload.Nickname))

	return err
}
