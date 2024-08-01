package bot

import (
	"bytes"
	"context"
	"fmt"
	"gopkg.in/telebot.v3/middleware"
	"image/jpeg"
	"log/slog"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"

	"github.com/WildEgor/pi-storyteller/internal/configs"
)

var _ Bot = (*TelegramBotAdapter)(nil)

// TelegramBotAdapter wrapper around telegram api
type TelegramBotAdapter struct {
	bot *tele.Bot

	callbacks   map[string]btnCallbackHandler
	defaultLang string
}

// NewTelegramBot ...
func NewTelegramBot(config *configs.TelegramBotConfig) *TelegramBotAdapter {
	pref := tele.Settings{
		Token:  config.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		slog.Error("bot init fail", slog.Any("err", err))
		panic("")
	}

	return &TelegramBotAdapter{
		bot:         bot,
		callbacks:   make(map[string]btnCallbackHandler),
		defaultLang: config.Language,
	}
}

// Start ...
func (t *TelegramBotAdapter) Start() {
	t.bot.Start()
}

// Stop ...
func (t *TelegramBotAdapter) Stop() {
	t.bot.Stop()
}

// SendMsg ...
func (t *TelegramBotAdapter) SendMsg(ctx context.Context, to *MessageRecipient, msg string) (mid int, err error) {
	m, err := t.bot.Send(to, msg)
	if err != nil {
		return 0, err
	}

	return m.ID, err
}

// EditMsg ...
func (t *TelegramBotAdapter) EditMsg(ctx context.Context, to *MessageRecipient, msg string) (mid int, err error) {
	m, err := t.bot.Edit(to, msg)
	if err != nil {
		return 0, err
	}

	return m.ID, err
}

// DeleteMsg ...
func (t *TelegramBotAdapter) DeleteMsg(ctx context.Context, to *MessageRecipient) error {
	return t.bot.Delete(to)
}

// SendStory ...
func (t *TelegramBotAdapter) SendStory(ctx context.Context, to *MessageRecipient, slides []StorySlide) error {
	files := make(tele.Album, 0, len(slides))

	sb := strings.Builder{}

	for i, v := range slides {
		buf := new(bytes.Buffer)
		//nolint
		_ = jpeg.Encode(buf, v.Image, nil)

		photo := &tele.Photo{
			File:    tele.FromReader(bytes.NewReader(buf.Bytes())),
			Caption: v.Desc,
		}

		if i == 0 {
			photo.Caption = fmt.Sprintf("[%s] %s", v.Style, v.Desc)
			if len(v.Style) == 0 {
				photo.Caption = v.Desc
			}
		}

		files = append(files, photo)

		//nolint
		sb.WriteString(v.Desc)
	}

	_, err := t.bot.SendAlbum(to, files)
	if err != nil {
		return err
	}

	return err
}

// HandleCommand ...
func (t *TelegramBotAdapter) HandleCommand(ctx context.Context, command string, fn commandHandler) {
	t.bot.Handle(command, func(c tele.Context) error {

		lang := c.Sender().LanguageCode
		if t.defaultLang != "" {
			lang = t.defaultLang
		}

		return fn(&CommandData{
			Nickname:  c.Sender().Username,
			MessageID: c.Message().ID,
			ChatID:    c.Chat().ID,
			Payload:   c.Message().Payload,
			Lang:      lang,
		})
	}, middleware.Recover())
}

// RegisterBtnCallback ...
func (t *TelegramBotAdapter) RegisterBtnCallback(ctx context.Context, name string, fn btnCallbackHandler) {
	t.callbacks[name] = fn
}
