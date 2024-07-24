package bot

import (
	"bytes"
	"context"
	"image/jpeg"
	"log/slog"
	"time"

	"github.com/WildEgor/pi-storyteller/internal/configs"
	tele "gopkg.in/telebot.v3"
)

var _ Bot = (*TelegramBot)(nil)

// TelegramBot wrapper around telegram api
type TelegramBot struct {
	bot *tele.Bot
}

// NewTelegramBot ...
func NewTelegramBot(config *configs.TelegramBotConfig) *TelegramBot {
	pref := tele.Settings{
		Token:  config.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		slog.Error("bot init fail", slog.Any("err", err))
		panic("")
	}

	return &TelegramBot{
		bot,
	}
}

// Start ...
func (t *TelegramBot) Start() {
	t.bot.Start()
}

// Stop ...
func (t *TelegramBot) Stop() {
	t.bot.Stop()
}

// SendMsg ...
func (t *TelegramBot) SendMsg(ctx context.Context, to *MessageRecipient, msg string) (mid int, err error) {
	m, err := t.bot.Send(to, msg)
	if err != nil {
		return 0, err
	}

	return m.ID, err
}

// EditMsg ...
func (t *TelegramBot) EditMsg(ctx context.Context, to *MessageRecipient, msg string) (mid int, err error) {
	m, err := t.bot.Edit(to, msg)
	if err != nil {
		return 0, err
	}

	return m.ID, err
}

// DeleteMsg ...
func (t *TelegramBot) DeleteMsg(ctx context.Context, to *MessageRecipient) error {
	return t.bot.Delete(to)
}

// SendStory ...
func (t *TelegramBot) SendStory(ctx context.Context, to *MessageRecipient, slides []StorySlide) error {
	files := make(tele.Album, 0, len(slides))

	for _, v := range slides {
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, v.Image, nil)

		photo := tele.Photo{
			File:    tele.FromReader(bytes.NewReader(buf.Bytes())),
			Caption: v.Desc,
		}
		files = append(files, &photo)
	}

	_, err := t.bot.SendAlbum(to, files)
	if err != nil {
		return err
	}

	return nil
}

// HandleCommand ...
func (t *TelegramBot) HandleCommand(ctx context.Context, command string, fn func(data *CommandData) error) {
	t.bot.Handle(command, func(c tele.Context) error {
		return fn(&CommandData{
			MessageID: c.Message().ID,
			ChatID:    c.Chat().ID,
			Payload:   c.Message().Payload,
		})
	})
}
