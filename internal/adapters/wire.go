package adapters

import (
	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	"github.com/WildEgor/pi-storyteller/internal/adapters/imaginator"
	"github.com/WildEgor/pi-storyteller/internal/adapters/textor"
	"github.com/google/wire"
)

// Set contains "adapters" to 3th party systems
var Set = wire.NewSet(
	imaginator.NewKandinskyAdapter,
	wire.Bind(new(imaginator.Imagininator), new(*imaginator.KandinskyAdapter)),
	textor.NewOllamaAdapter,
	wire.Bind(new(textor.ITextor), new(*textor.OllamaAdapter)),
	bot.NewTelegramBot,
	wire.Bind(new(bot.IBot), new(*bot.TelegramBot)),
	wire.Bind(new(bot.IBotRegistry), new(*bot.TelegramBot)),
)
