// Package tg_generate_handler is responsible for generating images
package tg_generate_handler

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sort"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	"github.com/WildEgor/pi-storyteller/internal/adapters/imaginator"
	"github.com/WildEgor/pi-storyteller/internal/configs"
	"github.com/WildEgor/pi-storyteller/internal/services/dispatcher"
	"github.com/WildEgor/pi-storyteller/internal/services/prompter"
	"github.com/WildEgor/pi-storyteller/internal/services/templater"
)

// GenerateHandler ...
type GenerateHandler struct {
	appConfig     *configs.AppConfig
	jobDispatcher dispatcher.Dispatcher
	imgGenerator  imaginator.Imagininator
	template      templater.Templater
	prompt        prompter.Prompter
	tgBot         bot.Bot
}

// NewGenerateHandler ...
func NewGenerateHandler(
	appConfig *configs.AppConfig,
	jobDispatcher dispatcher.Dispatcher,
	imgGenerator imaginator.Imagininator,
	template templater.Templater,
	prompt prompter.Prompter,
	tgBot bot.Bot,
) *GenerateHandler {
	return &GenerateHandler{
		appConfig,
		jobDispatcher,
		imgGenerator,
		template,
		prompt,
		tgBot,
	}
}

// Handle ...
func (h *GenerateHandler) Handle(ctx context.Context, payload *GenerateCommandDTO) error {
	chat := &bot.MessageRecipient{
		ID: payload.ChatID,
	}

	if err := payload.Validate(); err != nil {
		_, err := h.tgBot.SendMsg(ctx, chat, err.Error())
		return err
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
			_, err := h.tgBot.SendMsg(ctx, chat, "Please, wait! Too many request from you :)")
			return err
		}
	}

	slog.Info("new generate request", slog.Any("nickname", payload.Nickname), slog.Any("prompt", payload.Prompt))
	mid, err := h.tgBot.SendMsg(ctx, chat, "Start process... Please, wait!")
	chat.MessageID = strconv.Itoa(mid)

	id, err := h.jobDispatcher.Dispatch(func(jobCtx dispatcher.JobCtx) error {
		tCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		prompted := h.prompt.GetPredefinedRandomStyleStory(payload.Prompt, true)

		prompts := make([]string, 0, len(prompted))
		for _, conv := range prompted {
			prompts = append(prompts, conv.Prompt)
		}

		results := h.imgGenerator.GenerateImages(tCtx, prompts)

		errg := errgroup.Group{}
		errg.Go(func() error {
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
		})

		return errg.Wait()
	}, opts)

	slog.Info("dispatch task", slog.Any("uuid", id))

	if err != nil {
		//nolint
		_, _ = h.tgBot.SendMsg(ctx, chat, "Something goes wrong. Try later!")
	}

	return err
}
