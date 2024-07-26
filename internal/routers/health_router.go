// Package for health monitoring

package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/limiter"

	"log/slog"

	"github.com/WildEgor/pi-storyteller/internal/handlers/http/http_health_check_handler"
)

var _ IRouter[*fiber.App] = (*HealthRouter)(nil)

// HealthRouter router

type HealthRouter struct {
	hch *http_health_check_handler.HealthCheckHandler
}

// NewHealthRouter creates router

func NewHealthRouter(

	hch *http_health_check_handler.HealthCheckHandler,

) *HealthRouter {

	return &HealthRouter{

		hch,
	}

}

// Setup router

func (r *HealthRouter) Setup(app *fiber.App) {

	api := app.Group("/api", limiter.New(limiter.Config{

		Max: 10,

		SkipSuccessfulRequests: true,
	}))

	v1 := api.Group("/v1")

	v1.Get("/livez", healthcheck.NewHealthChecker(healthcheck.Config{

		Probe: func(ctx fiber.Ctx) bool {

			if err := r.hch.Handle(ctx); err != nil {

				slog.Error("error not healthy")

				return false

			}

			slog.Debug("is healthy")

			return true

		},
	}))

	v1.Get("/readyz", healthcheck.NewHealthChecker(healthcheck.Config{

		Probe: func(ctx fiber.Ctx) bool {

			return true

		},
	}))

}
