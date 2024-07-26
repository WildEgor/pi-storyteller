// Package routers route http/tg. Acts like "controllers" for routing http or etc.
package routers

import (
	"github.com/WildEgor/pi-storyteller/internal/handlers"
	"github.com/google/wire"
)

// Set ...
var Set = wire.NewSet(
	handlers.Set,
	NewHealthRouter,
	NewMetricsRouter,
	NewImageRouter,
)
