package pkg

import (
	"context"
	"fmt"
	"github.com/WildEgor/pi-storyteller/internal/services/dispatcher"
	"log/slog"
	"time"

	"github.com/WildEgor/pi-storyteller/internal/adapters"
	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	"github.com/WildEgor/pi-storyteller/internal/configs"
	eh "github.com/WildEgor/pi-storyteller/internal/handlers/http/http_error_handler"
	"github.com/WildEgor/pi-storyteller/internal/routers"
	"github.com/WildEgor/pi-storyteller/internal/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/google/wire"
)

// AppSet link main app deps
var AppSet = wire.NewSet(
	NewApp,
	configs.Set,
	routers.Set,
	adapters.Set,
	services.Set,
)

// Server represents the main server configuration.
type Server struct {
	App          *fiber.App
	Bot          bot.IBot
	Dispatcher   *dispatcher.Dispatcher
	Configurator *configs.Configurator
	AppConfig    *configs.AppConfig
}

// Run start service with deps
func (srv *Server) Run(ctx context.Context) {
	slog.Info("server is listening")

	srv.Dispatcher.Start()

	go func() {
		srv.Bot.Start()
		slog.Info("bot is listening")
	}()

	if err := srv.App.Listen(fmt.Sprintf(":%s", srv.AppConfig.HTTPPort), fiber.ListenConfig{
		DisableStartupMessage: false,
		EnablePrintRoutes:     false,
		OnShutdownSuccess: func() {
			slog.Debug("success shutdown service")
		},
	}); err != nil {
		slog.Error("unable to start server")
	}
}

// Shutdown graceful shutdown
func (srv *Server) Shutdown(ctx context.Context) {
	slog.Info("shutdown service")

	if err := srv.App.Shutdown(); err != nil {
		slog.Error("unable to shutdown server")
	}
}

// NewApp init app
func NewApp(
	ac *configs.AppConfig,
	lc *configs.LoggerConfig, // init logger
	c *configs.Configurator,
	eh *eh.ErrorsHandler,
	pbr *routers.HealthRouter,
	tr *routers.TelegramhRouter,
	bot bot.IBot,
	dispatcher *dispatcher.Dispatcher,
) *Server {
	app := fiber.New(fiber.Config{
		AppName:      ac.Name,
		ErrorHandler: eh.Handle,
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
	tr.Setup(app)

	return &Server{
		App:          app,
		Dispatcher:   dispatcher,
		Bot:          bot,
		AppConfig:    ac,
		Configurator: c,
	}
}
