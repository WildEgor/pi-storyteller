package tg_generate_handler

import (
	"context"
	"log/slog"
	"strings"

	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	"github.com/WildEgor/pi-storyteller/internal/adapters/imaginator"
	"github.com/WildEgor/pi-storyteller/internal/services/dispatcher"
)

type GenerateHandler struct {
	dispatcher *dispatcher.Dispatcher
	imaginator imaginator.Imagininator
	bot        bot.IBot
}

func NewGenerateHandler(dispatcher *dispatcher.Dispatcher, imaginator imaginator.Imagininator, bot bot.IBot) *GenerateHandler {
	return &GenerateHandler{
		dispatcher,
		imaginator,
		bot,
	}
}

func (h *GenerateHandler) Handle(ctx context.Context, payload *GeneratePayload) error {
	promts := strings.Split(payload.Prompt, ".")
	recp := &bot.MessageRecipient{
		ID: payload.ChatID,
	}

	if len(promts) > 5 {
		return h.bot.SendMsg(recp, "Ooops, limit reached")
	}

	err := h.dispatcher.Dispatch(func() {
		p := make(chan string, len(promts))
		r := make(chan imaginator.ImageGenerationResult)

		func() {
			defer close(p)
			for _, s := range promts {
				p <- s
			}
		}()

		h.imaginator.GenerateImages(ctx, p, r)

		images := make([]bot.StorySlide, 0)
		for v := range r {
			images = append(images, bot.StorySlide{
				Image: v.Image,
				Desc:  v.Prompt,
			})
		}

		err := h.bot.SendSlices(&bot.MessageRecipient{
			ID: payload.ChatID,
		}, images)
		if err != nil {
			slog.Error("error generating", slog.Any("err", err))
		}
	})

	return err
}
