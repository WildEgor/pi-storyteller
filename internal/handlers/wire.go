package handlers

import (
	"github.com/google/wire"

	eh "github.com/WildEgor/pi-storyteller/internal/handlers/http/http_error_handler"
	hch "github.com/WildEgor/pi-storyteller/internal/handlers/http/http_health_check_handler"
	tgh "github.com/WildEgor/pi-storyteller/internal/handlers/tg/generate"
	tsh "github.com/WildEgor/pi-storyteller/internal/handlers/tg/start"
)

// Set contains http/amqp/etc handlers (acts like facades)
var Set = wire.NewSet(
	eh.NewErrorsHandler,
	hch.NewHealthCheckHandler,
	tgh.NewGenerateHandler,
	tsh.NewStartHandler,
)
