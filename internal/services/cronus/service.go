package cronus

import (
	"context"
	"log/slog"
	"sort"
	"time"

	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	"github.com/WildEgor/pi-storyteller/internal/adapters/imaginator"
	"github.com/WildEgor/pi-storyteller/internal/configs"
	"github.com/WildEgor/pi-storyteller/internal/services/dispatcher"
	"github.com/WildEgor/pi-storyteller/internal/services/prompter"
)

const defaultTimezone = "Asia/Almaty"

// Service ...
type Service struct {
	dptchr       dispatcher.Dispatcher
	prompt       prompter.Prompter
	tgBot        bot.Bot
	imgGenerator imaginator.Imagininator

	tgConfig *configs.TelegramBotConfig

	crons []*dispatcher.DispatchCron
	loc   *time.Location
}

// New ...
func New(
	dptchr dispatcher.Dispatcher,
	prompt prompter.Prompter,
	tgBot bot.Bot,
	imgGenerator imaginator.Imagininator,
	tgConfig *configs.TelegramBotConfig,
) *Service {
	// TODO: specify using app config
	//nolint
	loc, _ := time.LoadLocation(defaultTimezone)

	return &Service{
		dptchr,
		prompt,
		tgBot,
		imgGenerator,
		tgConfig,
		make([]*dispatcher.DispatchCron, 0, 1),
		loc,
	}
}

// Start ...
func (s *Service) Start() {
	//nolint
	cron, _ := s.dptchr.DispatchCron(func(ctx dispatcher.JobCtx) error {
		slog.Info("SUCCESS news cron started")

		news, err := s.prompt.GetRandomNews()
		if err != nil {
			slog.Error("FAIL get news", slog.Any("err", err))
			return err
		}

		modifiedNew := s.prompt.GetPredefinedRandomStyleStory(news, false)

		tCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		slog.Debug("find news", slog.Any("value", modifiedNew))
		results := s.imgGenerator.GenerateImages(tCtx, []string{modifiedNew[0].Prompt})

		images := make([]bot.StorySlide, 0, len(results))
		for v := range results {
			images = append(images, bot.StorySlide{
				ID:    v.ID,
				Image: v.Image,
				Desc:  modifiedNew[0].Original,
			})
		}

		sort.Slice(images, func(i, j int) bool { return images[i].ID < images[j].ID })

		sErr := s.tgBot.SendStory(ctx, &bot.MessageRecipient{
			ID: s.tgConfig.ChatID,
		}, images)
		if sErr != nil {
			slog.Error("FAIL cron job", slog.Any("err", err))
			return sErr
		}

		slog.Info("SUCCESS news cron ended")

		return nil
	}, "0 * * * *", s.loc)

	s.crons = append(s.crons, cron)

	slog.Info("SUCCESS cron initialized")
}

// Stop ...
func (s *Service) Stop() {
	for _, cron := range s.crons {
		cron.Stop()
	}
}
