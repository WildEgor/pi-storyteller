// Package adapters contains "adapters" to 3rd party systems
package adapters

import (
	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	"github.com/WildEgor/pi-storyteller/internal/adapters/imaginator"
	"github.com/WildEgor/pi-storyteller/internal/adapters/textor"
	"github.com/WildEgor/pi-storyteller/pkg/kandinsky"
	"github.com/google/wire"
)

// Set ...
var Set = wire.NewSet(
	imaginator.NewKandinskyDummyClientProvider,
	// imaginator.NewKandinskyClientProvider,
	wire.Bind(new(kandinsky.Client), new(*imaginator.KandinskyClientProvider)),
	imaginator.NewKandinskyAdapter,
	wire.Bind(new(imaginator.Imagininator), new(*imaginator.KandinskyAdapter)),
	textor.NewOllamaAdapter,
	wire.Bind(new(textor.Textor), new(*textor.OllamaAdapter)),
	bot.NewTelegramBot,
	wire.Bind(new(bot.Bot), new(*bot.TelegramBot)),
	wire.Bind(new(bot.Registry), new(*bot.TelegramBot)),
)
