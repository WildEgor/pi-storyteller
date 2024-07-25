// Package route telegram commands
package routers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	tg_generate_handler "github.com/WildEgor/pi-storyteller/internal/handlers/tg/generate"
	tg_start_handler "github.com/WildEgor/pi-storyteller/internal/handlers/tg/start"
)

var _ IRouter[*fiber.App] = (*TelegramRouter)(nil)

// TelegramRouter router
type TelegramRouter struct {
	registry bot.Registry
	gh       *tg_generate_handler.GenerateHandler
	sh       *tg_start_handler.StartHandler
}

// NewImageRouter creates router
func NewImageRouter(
	registry bot.Registry,
	gh *tg_generate_handler.GenerateHandler,
	sh *tg_start_handler.StartHandler,
) *TelegramRouter {
	return &TelegramRouter{
		registry: registry,
		gh:       gh,
		sh:       sh,
	}
}

// Setup router
func (r *TelegramRouter) Setup(app *fiber.App) {
	r.registry.HandleCommand(context.TODO(), "/generate", func(data *bot.CommandData) error {
		return r.gh.Handle(context.TODO(), &tg_generate_handler.GenerateCommandDTO{
			Nickname:  data.Nickname,
			ChatID:    strconv.Itoa(int(data.ChatID)),
			MessageID: strconv.Itoa(int(data.MessageID)),
			Prompt:    data.Payload,
		})
	})

	r.registry.HandleCommand(context.TODO(), "/start", func(data *bot.CommandData) error {
		return r.sh.Handle(context.TODO(), &tg_start_handler.StartPayload{
			Nickname: data.Nickname,
			ChatID:   strconv.Itoa(int(data.ChatID)),
		})
	})
}
