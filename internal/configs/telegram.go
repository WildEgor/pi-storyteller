package configs

import (
	"log/slog"

	"github.com/spf13/viper"
)

// TelegramBotConfig holds telegram bot configuration
type TelegramBotConfig struct {
	Token string `mapstructure:"token"`
}

// NewTelegramBotConfig creates telegram bot config
func NewTelegramBotConfig(c *Configurator) *TelegramBotConfig {
	cfg := &TelegramBotConfig{}

	if err := viper.UnmarshalKey("telegram", &cfg); err != nil {
		panic("telegram config parse error")
	}

	slog.Info("config", slog.Any("value", cfg))

	return cfg
}
