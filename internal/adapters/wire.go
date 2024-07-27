// Package adapters contains "adapters" to 3rd party systems
package adapters

import (
	"github.com/google/wire"

	"github.com/WildEgor/pi-storyteller/internal/adapters/bot"
	"github.com/WildEgor/pi-storyteller/internal/adapters/imaginator"
	"github.com/WildEgor/pi-storyteller/internal/adapters/monitor"
	"github.com/WildEgor/pi-storyteller/internal/adapters/textor"
	"github.com/WildEgor/pi-storyteller/pkg/chatgpt"
	"github.com/WildEgor/pi-storyteller/pkg/kandinsky"
)

// Set ...
var Set = wire.NewSet(
	imaginator.NewKandinskyClientProvider,
	wire.Bind(new(kandinsky.Client), new(*imaginator.KandinskyClientProvider)),
	imaginator.NewKandinskyAdapter,
	wire.Bind(new(imaginator.Imagininator), new(*imaginator.KandinskyAdapter)),
	textor.NewChatGPTClientProvider,
	wire.Bind(new(chatgpt.Client), new(*textor.ChatGPTClientProvider)),
	textor.NewChatGPTAdapter,
	wire.Bind(new(textor.Textor), new(*textor.ChatGPTAdapter)),
	bot.NewTelegramBot,
	wire.Bind(new(bot.Bot), new(*bot.TelegramBotAdapter)),
	wire.Bind(new(bot.Registry), new(*bot.TelegramBotAdapter)),
	monitor.NewPromMetrics,
	wire.Bind(new(monitor.Monitor), new(*monitor.PromMetrics)),
)
