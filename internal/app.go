// Package app link main app deps
package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/google/wire"

	"github.com/WildEgor/pi-storyteller/internal/adapters"
	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	"github.com/WildEgor/pi-storyteller/internal/configs"
	eh "github.com/WildEgor/pi-storyteller/internal/handlers/http/http_error_handler"
	"github.com/WildEgor/pi-storyteller/internal/routers"
	"github.com/WildEgor/pi-storyteller/internal/services"
	"github.com/WildEgor/pi-storyteller/internal/services/cronus"
	"github.com/WildEgor/pi-storyteller/internal/services/dispatcher"
)

// Set ...
var Set = wire.NewSet(
	NewApp,
	configs.Set,
	routers.Set,
	adapters.Set,
	services.Set,
)

// App represents the main server configuration.
type App struct {
	App          *fiber.App
	Cron         cronus.Cronus
	Bot          bot.Bot
	Dispatcher   dispatcher.Dispatcher
	Configurator *configs.Configurator
	AppConfig    *configs.AppConfig
}

// Run start service with deps
func (srv *App) Run(_ context.Context) {
	go func() {
		slog.Info("dispatcher is listening")
		srv.Dispatcher.Start()
		slog.Info("cron is listening")
		srv.Cron.Start()
		slog.Info("bot is listening")
		srv.Bot.Start()
		// blocked
	}()

	slog.Info("server is listening")

	if err := srv.App.Listen(fmt.Sprintf(":%s", srv.AppConfig.HTTPPort), fiber.ListenConfig{
		DisableStartupMessage: false,
		EnablePrintRoutes:     false,
		OnShutdownSuccess: func() {
			slog.Info("success shutdown service")
		},
	}); err != nil {
		slog.Error("unable to start server", slog.Any("err", err))
	}
}

// Shutdown graceful shutdown
func (srv *App) Shutdown(_ context.Context) {
	slog.Info("shutdown service")
	srv.Bot.Stop()
	srv.Dispatcher.Stop()
	if err := srv.App.Shutdown(); err != nil {
		slog.Error("unable to shutdown server", slog.Any("err", err))
	}
}

// NewApp init app
func NewApp(
	c *configs.Configurator,
	ac *configs.AppConfig,
	lc *configs.LoggerConfig,
	erh *eh.ErrorsHandler,
	pbr *routers.HealthRouter,
	mr *routers.MetricsRouter,
	tr *routers.TelegramRouter,
	b bot.Bot,
	dptchr dispatcher.Dispatcher,
	crn cronus.Cronus,
) *App {
	app := fiber.New(fiber.Config{
		AppName:      ac.Name,
		ErrorHandler: erh.Handle,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS"},
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin"},
	}))

	app.Use(recover.New())

	pbr.Setup(app)
	mr.Setup(app)
	tr.Setup(app)

	return &App{
		App:          app,
		Cron:         crn,
		Dispatcher:   dptchr,
		Bot:          b,
		AppConfig:    ac,
		Configurator: c,
	}
}
