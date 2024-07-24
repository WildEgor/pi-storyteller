package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	server "github.com/WildEgor/pi-storyteller/internal"
)

// @title			[Service name here] Swagger Doc
// @version			1.0
// @description		[Service name here]
// @termsOfService	/
// @contact.name	mail
// @contact.url		/
// @contact.email	kartashov_egor96@mail.ru
// @license.name	MIT
// @license.url		http://www.apache.org/licenses/MIT.html
// @host			localhost:8888
// @BasePath		/
// @schemes			http
func main() {
	// Catch terminate signals
	ctx, done := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer done()

	srv, err := server.NewServer()
	if err != nil {
		slog.Error("srv fail", slog.Any("err", err))
		panic("")
	}

	srv.Run(ctx)

	<-ctx.Done()

	// Wait before shutdown
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer func() {
		cancel()
	}()

	srv.Shutdown(ctx)
}
