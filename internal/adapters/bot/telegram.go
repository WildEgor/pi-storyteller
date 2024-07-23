package bot

import (
	"bytes"
	"image/jpeg"
	"log"
	"time"

	"github.com/WildEgor/pi-storyteller/internal/configs"
	tele "gopkg.in/telebot.v3"
)

var _ IBot = (*TelegramBot)(nil)

type TelegramBot struct {
	bot *tele.Bot
}

func NewTelegramBot(config *configs.TelegramBotConfig) *TelegramBot {
	pref := tele.Settings{
		Token:  config.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	return &TelegramBot{
		bot,
	}
}

func (t *TelegramBot) Start() {
	t.bot.Start()
}

func (t *TelegramBot) Stop() {
	t.bot.Stop()
}

// SendMsg ...
func (t *TelegramBot) SendMsg(to *MessageRecipient, msg string) (mid int, err error) {
	m, err := t.bot.Send(to, msg)
	if err != nil {
		return 0, err
	}

	return m.ID, err
}

func (t *TelegramBot) EditMsg(to *MessageRecipient, msg string) (mid int, err error) {
	m, err := t.bot.Edit(to, msg)
	if err != nil {
		return 0, err
	}

	return m.ID, err
}

// SendStory ...
func (t *TelegramBot) SendStory(to *MessageRecipient, slides []StorySlide) error {
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

func (t *TelegramBot) HandleCommand(command string, fn func(data *CommandData) error) {
	t.bot.Handle(command, func(c tele.Context) error {
		return fn(&CommandData{
			ChatID:  c.Chat().ID,
			Payload: c.Message().Payload,
		})
	})
}
