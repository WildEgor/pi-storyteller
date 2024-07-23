package routers

import (
	"github.com/WildEgor/pi-storyteller/internal/handlers"
	"github.com/google/wire"
)

// Set acts like "controllers" for routing http or etc.
var Set = wire.NewSet(
	handlers.Set,
	NewHealthRouter,
	NewImageRouter,
)
