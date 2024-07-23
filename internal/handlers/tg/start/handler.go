package tg_start_handler

import (
	"context"

	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
)

type StartHandler struct {
	bot bot.IBot
}

func NewStartHandler(bot bot.IBot) *StartHandler {
	return &StartHandler{
		bot,
	}
}

func (h *StartHandler) Handle(ctx context.Context, payload *StartPayload) error {

	h.bot.SendMsg(&bot.MessageRecipient{
		ID: payload.ChatID,
	}, "HI!")

	return nil
}
