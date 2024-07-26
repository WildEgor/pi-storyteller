// Package routers route http/tg. Acts like "controllers" for routing http or etc.
package routers

import (
	"github.com/google/wire"

	"github.com/WildEgor/pi-storyteller/internal/handlers"
)

// Set ...
var Set = wire.NewSet(
	handlers.Set,
	NewHealthRouter,
	NewMetricsRouter,
	NewImageRouter,
)
