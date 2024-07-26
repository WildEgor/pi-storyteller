package http_metrics_handler

import (
	"github.com/WildEgor/pi-storyteller/internal/adapters/monitor"
	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type MetricsHandler struct {
	reg *monitor.PromMetricsRegistry
}

func NewMetricsHandler(reg *monitor.PromMetricsRegistry) *MetricsHandler {
	return &MetricsHandler{
		reg: reg,
	}
}

// Handle MetricsHandler godoc
//
//	@Summary		Metrics service
//	@Description	Metrics service
//	@Tags			Metrics Controller
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/metrics [get]
func (h *MetricsHandler) Handle(ctx fiber.Ctx) error {
	ph := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())
	ph(ctx.Context())
	return nil
}
