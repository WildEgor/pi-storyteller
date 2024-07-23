package services

import (
	"github.com/WildEgor/pi-storyteller/internal/services/dispatcher"
	"github.com/google/wire"
)

// Set contains "services"
var Set = wire.NewSet(
	dispatcher.NewDispatcher,
)
