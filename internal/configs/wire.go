// Package configs contains any kind of configs
package configs

import "github.com/google/wire"

// Set ...
var Set = wire.NewSet(
	NewConfigurator,
	NewAppConfig,
	NewLoggerConfig,
	NewKandinskyConfig,
	NewChatGPTConfig,
	NewTelegramBotConfig,
	NewMetricsConfig,
)
