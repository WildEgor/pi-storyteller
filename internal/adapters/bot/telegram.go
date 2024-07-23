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

// SendMsg implements IBot.
func (t *TelegramBot) SendMsg(to *MessageRecipient, msg string) error {
	_, err := t.bot.Send(to, msg)

	return err
}

// SendSlices implements IBot.
func (t *TelegramBot) SendSlices(to *MessageRecipient, slides []StorySlide) error {
	files := make(tele.Album, 0, len(slides))

	for _, v := range slides {
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, v.Image, nil)

		photo := tele.Photo{File: tele.FromReader(bytes.NewReader(buf.Bytes()))}
		files = append(files, &photo)
	}

	_, err := t.bot.SendAlbum(to, files)
	if err != nil {
		return err
	}

	return nil
}

func (t *TelegramBot) HandleCommand(command string, fn func(data *TelegramCommandData)) error {
	t.bot.Handle(command, func(c tele.Context) error {
		fn(&TelegramCommandData{
			ChatID:  c.Chat().ID,
			Payload: c.Message().Payload,
		})

		return nil
	})

	return nil
}
