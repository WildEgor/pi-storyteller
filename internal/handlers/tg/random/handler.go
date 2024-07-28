// Package tg_random_handler is responsible for random generating images
package tg_random_handler

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sort"
	"strconv"
	"time"

	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	"github.com/WildEgor/pi-storyteller/internal/adapters/imaginator"
	"github.com/WildEgor/pi-storyteller/internal/configs"
	"github.com/WildEgor/pi-storyteller/internal/services/dispatcher"
	"github.com/WildEgor/pi-storyteller/internal/services/prompter"
	"github.com/WildEgor/pi-storyteller/internal/services/templater"
)

// RandomHandler ...
type RandomHandler struct {
	appConfig     *configs.AppConfig
	jobDispatcher dispatcher.Dispatcher
	imgGenerator  imaginator.Imagininator
	template      templater.Templater
	prompt        prompter.Prompter
	tgBot         bot.Bot
}

// NewRandomHandler ...
func NewRandomHandler(
	appConfig *configs.AppConfig,
	jobDispatcher dispatcher.Dispatcher,
	imgGenerator imaginator.Imagininator,
	template templater.Templater,
	prompt prompter.Prompter,
	tgBot bot.Bot,
) *RandomHandler {
	return &RandomHandler{
		appConfig,
		jobDispatcher,
		imgGenerator,
		template,
		prompt,
		tgBot,
	}
}

// Handle ...
func (h *RandomHandler) Handle(ctx context.Context, payload *RandomCommandDTO) error {
	chat := &bot.MessageRecipient{
		ID: payload.ChatID,
	}

	opts := &dispatcher.JobOpts{
		OwnerID: payload.Nickname,
	}
	if slices.Contains(h.appConfig.PriorityList, payload.Nickname) {
		opts.Priority = dispatcher.HighPriority
	} else {
		count := h.jobDispatcher.CountActiveJobs(payload.Nickname)
		if count > 3 {
			slog.Warn(fmt.Sprintf("%s still wait. Has %d active jobs", payload.Nickname, count))
			_, err := h.tgBot.SendMsg(ctx, chat, "🤯🤯🤯")
			return err
		}
	}

	slog.Info("new random request", slog.Any("nickname", payload.Nickname), slog.Any("lang", payload.Lang))
	mid, err := h.tgBot.SendMsg(ctx, chat, "🤔")
	chat.MessageID = strconv.Itoa(mid)

	id, err := h.jobDispatcher.Dispatch(func(jobCtx dispatcher.JobCtx) error {
		tCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		prompted := h.prompt.GetRandomStory(payload.Lang)
		prompts := make([]string, 0, len(prompted))
		for _, conv := range prompted {
			prompts = append(prompts, conv.Prompt)
		}

		results := h.imgGenerator.GenerateImages(tCtx, prompts)

		//nolint
		_ = h.tgBot.DeleteMsg(ctx, chat)
		//nolint
		_ = h.tgBot.DeleteMsg(ctx, &bot.MessageRecipient{
			ID:        payload.ChatID,
			MessageID: payload.MessageID,
		})

		images := make([]bot.StorySlide, 0, len(results))
		for v := range results {
			images = append(images, bot.StorySlide{
				ID:    v.ID,
				Style: prompted[0].Style,
				Image: v.Image,
				Desc:  prompted[v.ID].Original,
			})
		}

		sort.Slice(images, func(i, j int) bool { return images[i].ID < images[j].ID })

		sErr := h.tgBot.SendStory(ctx, &bot.MessageRecipient{
			ID: payload.ChatID,
		}, images)
		if sErr != nil {
			slog.Error("error generating", slog.Any("err", err))
		}

		return sErr
	}, opts)

	slog.Info("dispatch task", slog.Any("uuid", id))

	if err != nil {
		//nolint
		_, _ = h.tgBot.SendMsg(ctx, chat, "Something goes wrong. Try later!")
	}

	return err
}
