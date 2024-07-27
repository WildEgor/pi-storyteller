// Package http_metrics_handler ...
package http_metrics_handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// MetricsHandler ...
type MetricsHandler struct {
}

// NewMetricsHandler ...
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
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
