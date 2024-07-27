// Package routers route telegram commands
package routers

import (
	"github.com/gofiber/fiber/v3"

	"context"
	"strconv"

	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	tg_generate_handler "github.com/WildEgor/pi-storyteller/internal/handlers/tg/generate"
	tg_random_handler "github.com/WildEgor/pi-storyteller/internal/handlers/tg/random"
	tg_start_handler "github.com/WildEgor/pi-storyteller/internal/handlers/tg/start"
)

var _ IRouter[*fiber.App] = (*TelegramRouter)(nil)

// TelegramRouter router
type TelegramRouter struct {
	br bot.Registry
	gh *tg_generate_handler.GenerateHandler
	sh *tg_start_handler.StartHandler
	rh *tg_random_handler.RandomHandler
}

// NewImageRouter creates router
func NewImageRouter(
	br bot.Registry,
	gh *tg_generate_handler.GenerateHandler,
	sh *tg_start_handler.StartHandler,
	rh *tg_random_handler.RandomHandler,
) *TelegramRouter {
	return &TelegramRouter{
		br: br,
		gh: gh,
		sh: sh,
		rh: rh,
	}
}

// Setup router
func (r *TelegramRouter) Setup(_ *fiber.App) {
	r.br.HandleCommand(context.TODO(), "/generate", func(data *bot.CommandData) error {
		return r.gh.Handle(context.TODO(), &tg_generate_handler.GenerateCommandDTO{
			Nickname:  data.Nickname,
			ChatID:    strconv.Itoa(int(data.ChatID)),
			MessageID: strconv.Itoa(data.MessageID),
			Prompt:    data.Payload,
		})
	})

	r.br.HandleCommand(context.TODO(), "/random", func(data *bot.CommandData) error {
		return r.rh.Handle(context.TODO(), &tg_random_handler.RandomCommandDTO{
			Nickname:  data.Nickname,
			ChatID:    strconv.Itoa(int(data.ChatID)),
			MessageID: strconv.Itoa(data.MessageID),
			Lang:      data.Lang,
		})
	})

	r.br.HandleCommand(context.TODO(), "/start", func(data *bot.CommandData) error {
		return r.sh.Handle(context.TODO(), &tg_start_handler.StartDTO{
			Nickname: data.Nickname,
			ChatID:   strconv.Itoa(int(data.ChatID)),
		})
	})
}
