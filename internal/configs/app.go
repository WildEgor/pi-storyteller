package configs

import (
	"log/slog"

	"github.com/spf13/viper"
)

// AppConfig holds the main app configurations
type AppConfig struct {
	Name     string `mapstructure:"name"`
	Mode     string `mapstructure:"mode"`
	HTTPPort string `mapstructure:"http_port"`
	changer  func(ac *AppConfig)
}

// NewAppConfig creates app config
func NewAppConfig() *AppConfig {
	cfg := &AppConfig{}

	if err := viper.UnmarshalKey("app", &cfg); err != nil {
		slog.Error("app parse error", slog.Any("err", err))
		panic("")
	}

	slog.Debug("app config", slog.Any("value", cfg))

	return cfg
}

// OnChanged register callback
func (ac *AppConfig) OnChanged(fn func(ac *AppConfig)) {
	ac.changer = fn
}

// IsProduction Check is application running in production mode
func (ac *AppConfig) IsProduction() bool {
	return ac.Mode == "production"
}

// IsDebug Check is application running in debug mode
func (ac *AppConfig) IsDebug() bool {
	return ac.Mode == "debug"
}
