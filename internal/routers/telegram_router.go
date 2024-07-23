package routers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	tg_generate_handler "github.com/WildEgor/pi-storyteller/internal/handlers/tg/generate"
	tg_start_handler "github.com/WildEgor/pi-storyteller/internal/handlers/tg/start"
)

var _ IRouter[*fiber.App] = (*HealthRouter)(nil)

// TelegramhRouter router
type TelegramhRouter struct {
	registry bot.ITelegramBotRegistry
	gh       *tg_generate_handler.GenerateHandler
	sh       *tg_start_handler.StartHandler
}

// NewHealthRouter creates router
func NewImageRouter(
	registry bot.ITelegramBotRegistry,
	gh *tg_generate_handler.GenerateHandler,
	sh *tg_start_handler.StartHandler,
) *TelegramhRouter {
	return &TelegramhRouter{
		registry: registry,
		gh:       gh,
		sh:       sh,
	}
}

// Setup router
func (r *TelegramhRouter) Setup(app *fiber.App) {
	r.registry.HandleCommand("/generate", func(data *bot.TelegramCommandData) {
		r.gh.Handle(context.TODO(), &tg_generate_handler.GeneratePayload{
			ChatID: strconv.Itoa(int(data.ChatID)),
			Prompt: data.Payload,
		})
	})

	r.registry.HandleCommand("/start", func(data *bot.TelegramCommandData) {
		r.sh.Handle(context.TODO(), &tg_start_handler.StartPayload{
			ChatID: strconv.Itoa(int(data.ChatID)),
		})
	})
}
