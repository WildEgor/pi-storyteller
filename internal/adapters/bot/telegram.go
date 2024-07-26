package bot

import (
	"bytes"
	"context"
	"image/jpeg"
	"log/slog"
	"strings"
	"time"

	"github.com/WildEgor/pi-storyteller/internal/configs"
	tele "gopkg.in/telebot.v3"
)

var _ Bot = (*TelegramBot)(nil)

// TelegramBot wrapper around telegram api
type TelegramBot struct {
	bot *tele.Bot

	callbacks map[string]btnCallbackHandler
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
		bot:       bot,
		callbacks: make(map[string]btnCallbackHandler),
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

	sb := strings.Builder{}

	for _, v := range slides {
		buf := new(bytes.Buffer)
		//nolint
		_ = jpeg.Encode(buf, v.Image, nil)

		files = append(files, &tele.Photo{
			File:    tele.FromReader(bytes.NewReader(buf.Bytes())),
			Caption: v.Desc,
		})

		sb.WriteString(v.Desc)
	}

	_, err := t.bot.SendAlbum(to, files)
	if err != nil {
		return err
	}

	selector := &tele.ReplyMarkup{}
	btnRegen := selector.Data("   🔁   ", "__generate_regenerate__", sb.String())
	btnWarn := selector.Data("   ⚠️   ", "__generate_warn__", sb.String())
	selector.Inline(
		selector.Row(btnRegen),
		selector.Row(btnWarn),
	)

	if fn, ok := t.callbacks[btnRegen.Unique]; ok {
		t.bot.Handle(&btnRegen, func(c tele.Context) error {
			slog.Debug("regen pressed")
			return fn(&BtnCallbackData{
				// TODO
			})
		})
	}

	if fn, ok := t.callbacks[btnWarn.Unique]; ok {
		t.bot.Handle(&btnWarn, func(c tele.Context) error {
			slog.Debug("warn pressed")
			return fn(&BtnCallbackData{
				// TODO
			})
		})
	}

	_, err = t.bot.Send(to, "~", &tele.SendOptions{
		ReplyMarkup: selector,
	})

	return err
}

// HandleCommand ...
func (t *TelegramBot) HandleCommand(ctx context.Context, command string, fn commandHandler) {
	t.bot.Handle(command, func(c tele.Context) error {
		return fn(&CommandData{
			Nickname:  c.Sender().Username,
			MessageID: c.Message().ID,
			ChatID:    c.Chat().ID,
			Payload:   c.Message().Payload,
		})
	})
}

// RegisterBtnCallback ...
func (t *TelegramBot) RegisterBtnCallback(ctx context.Context, name string, fn btnCallbackHandler) {
	t.callbacks[name] = fn
}
