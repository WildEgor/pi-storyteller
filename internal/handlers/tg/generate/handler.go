package tg_generate_handler

import (
	"context"
	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	"github.com/WildEgor/pi-storyteller/internal/adapters/imaginator"
	"github.com/WildEgor/pi-storyteller/internal/services/dispatcher"
	"github.com/WildEgor/pi-storyteller/internal/services/prompter"
	"github.com/WildEgor/pi-storyteller/internal/services/templater"
	"log/slog"
	"sort"
	"strconv"
	"time"
)

// GenerateHandler ...
type GenerateHandler struct {
	dispatcher *dispatcher.Dispatcher
	imaginator imaginator.Imagininator
	templater  *templater.Templater
	prompter   *prompter.Prompter
	bot        bot.Bot
}

// NewGenerateHandler ...
func NewGenerateHandler(
	dispatcher *dispatcher.Dispatcher,
	imaginator imaginator.Imagininator,
	templater *templater.Templater,
	prompter *prompter.Prompter,
	bot bot.Bot,
) *GenerateHandler {
	return &GenerateHandler{
		dispatcher,
		imaginator,
		templater,
		prompter,
		bot,
	}
}

// Handle ...
func (h *GenerateHandler) Handle(ctx context.Context, payload *GeneratePayload) error {
	chat := &bot.MessageRecipient{
		ID: payload.ChatID,
	}

	if err := payload.Validate(); err != nil {
		_, err := h.bot.SendMsg(ctx, chat, err.Error())
		return err
	}

	slog.Info("new request", slog.Any("nickname", payload.Nickname), slog.Any("prompt", payload.Prompt))

	// TODO: before new generate check if user queue already has await process
	// or in progress tasks (limit 5 per 30 min, except by whitelist)

	_, err := h.dispatcher.Dispatch(func() {
		mid, _ := h.bot.SendMsg(ctx, chat, "Start process... Please, wait!")
		chat.MessageID = strconv.Itoa(mid)

		tCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		prompted := h.prompter.Random(payload.Prompt)

		var prompts []string
		for _, conv := range prompted {
			prompts = append(prompts, conv.Prompt)
		}

		results := h.imaginator.GenerateImages(tCtx, prompts)

		go func() {
			h.bot.DeleteMsg(ctx, chat)
			h.bot.DeleteMsg(ctx, &bot.MessageRecipient{
				ID:        payload.ChatID,
				MessageID: payload.MessageID,
			})

			images := make([]bot.StorySlide, 0)
			for v := range results {
				images = append(images, bot.StorySlide{
					ID:    v.ID,
					Image: v.Image,
					Desc:  prompted[v.ID].Original,
				})
			}

			sort.Slice(images, func(i, j int) bool { return images[i].ID < images[j].ID })

			err := h.bot.SendStory(ctx, &bot.MessageRecipient{
				ID: payload.ChatID,
			}, images)
			if err != nil {
				slog.Error("error generating", slog.Any("err", err))
			}
		}()
	})

	return err
}
