package configs

import (
	"github.com/spf13/viper"

	"log/slog"
)

// TelegramBotConfig holds telegram bot configuration
type TelegramBotConfig struct {
	Token string `mapstructure:"token"`
}

// NewTelegramBotConfig creates telegram bot config
func NewTelegramBotConfig() *TelegramBotConfig {
	cfg := &TelegramBotConfig{}

	if err := viper.UnmarshalKey("telegram", &cfg); err != nil {
		slog.Error("telegram parse error", slog.Any("err", err))
		panic("")
	}

	slog.Debug("telegram config", slog.Any("value", cfg))

	return cfg
}
