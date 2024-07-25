// Package is responsible for generating images
package tg_generate_handler

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
	"golang.org/x/sync/errgroup"
)

// GenerateHandler ...
type GenerateHandler struct {
	appConfig  *configs.AppConfig
	dispatcher *dispatcher.Dispatcher
	imaginator imaginator.Imagininator
	templater  *templater.Templater
	prompter   *prompter.Prompter
	bot        bot.Bot
}

// NewGenerateHandler ...
func NewGenerateHandler(
	appConfig *configs.AppConfig,
	dispatcher *dispatcher.Dispatcher,
	imaginator imaginator.Imagininator,
	templater *templater.Templater,
	prompter *prompter.Prompter,
	bot bot.Bot,
) *GenerateHandler {
	return &GenerateHandler{
		appConfig,
		dispatcher,
		imaginator,
		templater,
		prompter,
		bot,
	}
}

// Handle ...
func (h *GenerateHandler) Handle(ctx context.Context, payload *GenerateCommandDTO) error {
	chat := &bot.MessageRecipient{
		ID: payload.ChatID,
	}

	if err := payload.Validate(); err != nil {
		_, err := h.bot.SendMsg(ctx, chat, err.Error())
		return err
	}

	opts := &dispatcher.JobOpts{
		OwnerID: payload.Nickname,
	}
	if slices.Contains(h.appConfig.PriorityList, payload.Nickname) {
		opts.Priority = dispatcher.HighPriority
	} else {
		count := h.dispatcher.CountActiveJobs(payload.Nickname)
		if count > 3 {
			slog.Warn(fmt.Sprintf("%s still wait. Has %d active jobs", payload.Nickname, count))
			_, err := h.bot.SendMsg(ctx, chat, "Please, wait! Too many request from you :)")
			return err
		}
	}

	slog.Info("new request", slog.Any("nickname", payload.Nickname), slog.Any("prompt", payload.Prompt))
	mid, err := h.bot.SendMsg(ctx, chat, "Start process... Please, wait!")
	chat.MessageID = strconv.Itoa(mid)

	id, err := h.dispatcher.Dispatch(func(jobCtx dispatcher.JobCtx) error {
		tCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		prompted := h.prompter.Random(payload.Prompt)

		var prompts []string
		for _, conv := range prompted {
			prompts = append(prompts, conv.Prompt)
		}

		results := h.imaginator.GenerateImages(tCtx, prompts)

		errg := errgroup.Group{}

		errg.Go(func() error {
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

			return err
		})

		err = errg.Wait()

		return err
	}, opts)

	slog.Info("dispatch task", slog.Any("uuid", id))

	return err
}
