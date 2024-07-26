package cronus

import (
	"context"
	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	"github.com/WildEgor/pi-storyteller/internal/adapters/imaginator"
	"github.com/WildEgor/pi-storyteller/internal/configs"
	"github.com/WildEgor/pi-storyteller/internal/services/dispatcher"
	"github.com/WildEgor/pi-storyteller/internal/services/prompter"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"sort"
	"time"
)

// Service ...
type Service struct {
	dptchr       *dispatcher.Dispatcher
	prompt       *prompter.Prompter
	tgBot        bot.Bot
	imgGenerator imaginator.Imagininator
	tgConfig     *configs.TelegramBotConfig

	crons []*dispatcher.DispatchCron
}

// New ...
func New(
	crons *dispatcher.Dispatcher,
	prompt *prompter.Prompter,
	tgBot bot.Bot,
	imgGenerator imaginator.Imagininator,
	tgConfig *configs.TelegramBotConfig,
) *Service {
	s := &Service{
		crons,
		prompt,
		tgBot,
		imgGenerator,
		tgConfig,
		nil,
	}

	s.crons = make([]*dispatcher.DispatchCron, 0, 1)
	return s
}

// Start ...
func (s *Service) Start() {
	loc, _ := time.LoadLocation("Asia/Almaty")

	cron, _ := s.dptchr.DispatchCron(func(ctx dispatcher.JobCtx) error {
		txt, err := s.prompt.GetRandomNews()
		if err != nil {
			slog.Error("fail get news", slog.Any("err", err))
			return err
		}

		tCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		slog.Debug("news: ", slog.Any("value", txt))
		results := s.imgGenerator.GenerateImages(tCtx, []string{txt})

		errg := errgroup.Group{}
		errg.Go(func() error {
			images := make([]bot.StorySlide, 0, len(results))
			for v := range results {
				images = append(images, bot.StorySlide{
					ID:    v.ID,
					Image: v.Image,
					Desc:  txt,
				})
			}

			sort.Slice(images, func(i, j int) bool { return images[i].ID < images[j].ID })

			sErr := s.tgBot.SendStory(ctx, &bot.MessageRecipient{
				ID: s.tgConfig.ChatID,
			}, images)
			if sErr != nil {
				slog.Error("error cron job", slog.Any("err", err))
				return sErr
			}

			return nil
		})

		return errg.Wait()

	}, "0 * * * * *", loc)

	slog.Info("init cron")

	s.crons = append(s.crons, cron)
}

// Stop ...
func (s *Service) Stop() {
	for _, cron := range s.crons {
		cron.Stop()
	}
}
