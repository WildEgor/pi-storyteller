// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	"github.com/WildEgor/pi-storyteller/internal/adapters/imaginator"
	"github.com/WildEgor/pi-storyteller/internal/adapters/monitor"
	"github.com/WildEgor/pi-storyteller/internal/configs"
	"github.com/WildEgor/pi-storyteller/internal/handlers/http/http_error_handler"
	"github.com/WildEgor/pi-storyteller/internal/handlers/http/http_health_check_handler"
	"github.com/WildEgor/pi-storyteller/internal/handlers/http/metrics"
	"github.com/WildEgor/pi-storyteller/internal/handlers/tg/generate"
	"github.com/WildEgor/pi-storyteller/internal/handlers/tg/start"
	"github.com/WildEgor/pi-storyteller/internal/routers"
	"github.com/WildEgor/pi-storyteller/internal/services/dispatcher"
	"github.com/WildEgor/pi-storyteller/internal/services/prompter"
	"github.com/WildEgor/pi-storyteller/internal/services/templater"
	"github.com/google/wire"
)

// Injectors from server.go:

// NewServer
func NewServer() (*App, error) {
	configurator := configs.NewConfigurator()
	appConfig := configs.NewAppConfig()
	loggerConfig := configs.NewLoggerConfig()
	errorsHandler := http_error_handler.NewErrorsHandler()
	healthCheckHandler := http_health_check_handler.NewHealthCheckHandler()
	healthRouter := routers.NewHealthRouter(healthCheckHandler)
	metricsConfig := configs.NewMetricsConfig()
	promMetricsRegistry := monitor.NewPromMetricsRegistry(metricsConfig)
	metricsHandler := http_metrics_handler.NewMetricsHandler(promMetricsRegistry)
	metricsRouter := routers.NewMetricsRouter(metricsHandler, metricsConfig)
	telegramBotConfig := configs.NewTelegramBotConfig()
	telegramBot := bot.NewTelegramBot(telegramBotConfig)
	promMetrics := monitor.NewPromMetrics(promMetricsRegistry, appConfig, metricsConfig)
	dispatcherDispatcher := dispatcher.NewDispatcher(promMetrics)
	kandinskyConfig := configs.NewKandinskyConfig()
	kandinskyClientProvider := imaginator.NewKandinskyDummyClientProvider(kandinskyConfig)
	kandinskyAdapter := imaginator.NewKandinskyAdapter(kandinskyClientProvider)
	templaterTemplater := templater.NewTemplateService(appConfig)
	prompterPrompter := prompter.New(appConfig)
	generateHandler := tg_generate_handler.NewGenerateHandler(appConfig, dispatcherDispatcher, kandinskyAdapter, templaterTemplater, prompterPrompter, telegramBot)
	startHandler := tg_start_handler.NewStartHandler(telegramBot)
	telegramRouter := routers.NewImageRouter(telegramBot, generateHandler, startHandler)
	app := NewApp(configurator, appConfig, loggerConfig, errorsHandler, healthRouter, metricsRouter, telegramRouter, telegramBot, dispatcherDispatcher)
	return app, nil
}

// server.go:

// ServerSet
var ServerSet = wire.NewSet(Set)
