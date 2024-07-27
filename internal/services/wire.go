// Package services contains "services"
package services

import (
	"github.com/WildEgor/pi-storyteller/internal/services/cronus"
	"github.com/WildEgor/pi-storyteller/internal/services/dispatcher"
	"github.com/WildEgor/pi-storyteller/internal/services/prompter"
	"github.com/WildEgor/pi-storyteller/internal/services/templater"
	"github.com/google/wire"
)

// Set ...
var Set = wire.NewSet(
	dispatcher.New,
	wire.Bind(new(dispatcher.Dispatcher), new(*dispatcher.Service)),
	templater.New,
	wire.Bind(new(templater.Templater), new(*templater.Service)),
	prompter.New,
	wire.Bind(new(prompter.Prompter), new(*prompter.Service)),
	cronus.New,
	wire.Bind(new(cronus.Cronus), new(*cronus.Service)),
)
