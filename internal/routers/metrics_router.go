// Package routers Package for health monitoring
package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"

	"github.com/WildEgor/pi-storyteller/internal/configs"
	http_metrics_handler "github.com/WildEgor/pi-storyteller/internal/handlers/http/metrics"
)

var _ IRouter[*fiber.App] = (*MetricsRouter)(nil)

// MetricsRouter router
type MetricsRouter struct {
	mh   *http_metrics_handler.MetricsHandler
	mcfg *configs.MetricsConfig
}

// NewMetricsRouter creates router
func NewMetricsRouter(
	mh *http_metrics_handler.MetricsHandler,
	mcfg *configs.MetricsConfig,
) *MetricsRouter {
	return &MetricsRouter{
		mh,
		mcfg,
	}
}

// Setup router
func (r *MetricsRouter) Setup(app *fiber.App) {
	api := app.Group("/api", limiter.New(limiter.Config{
		Max:                    10,
		SkipSuccessfulRequests: true,
	}))

	v1 := api.Group("/v1")

	if r.mcfg.Enabled {
		v1.Get("/metrics", r.mh.Handle)
	}
}
