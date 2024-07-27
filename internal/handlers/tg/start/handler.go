// Package tg_start_handler responsible to show instructions
package tg_start_handler

import (
	"context"
	"fmt"

	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
)

// StartHandler ...
type StartHandler struct {
	tgBot bot.Bot
}

// NewStartHandler ...
func NewStartHandler(tgBot bot.Bot) *StartHandler {
	return &StartHandler{
		tgBot,
	}
}

// Handle ...
func (h *StartHandler) Handle(ctx context.Context, payload *StartDTO) error {
	//nolint
	_, err := h.tgBot.SendMsg(ctx, &bot.MessageRecipient{
		ID: payload.ChatID,
	}, fmt.Sprintf("Hi, %s", payload.Nickname))

	return err
}
