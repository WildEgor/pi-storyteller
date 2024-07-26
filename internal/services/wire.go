// Package services contains "services"
package services

import (
	"github.com/google/wire"

	"github.com/WildEgor/pi-storyteller/internal/services/dispatcher"
	"github.com/WildEgor/pi-storyteller/internal/services/prompter"
	"github.com/WildEgor/pi-storyteller/internal/services/templater"
)

// Set ...
var Set = wire.NewSet(
	dispatcher.NewDispatcher,
	templater.NewTemplateService,
	prompter.New,
)
