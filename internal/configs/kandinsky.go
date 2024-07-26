package configs

import (
	"github.com/spf13/viper"

	"log/slog"

	"github.com/WildEgor/pi-storyteller/pkg/kandinsky"
)

// KandinskyConfig holds kandinsky api configuration
type KandinskyConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Key     string `mapstructure:"key"`
	Secret  string `mapstructure:"secret"`
	Debug   bool   `mapstructure:"debug"`
}

// NewKandinskyConfig creates kandinsky config
func NewKandinskyConfig() *KandinskyConfig {
	cfg := &KandinskyConfig{}

	if err := viper.UnmarshalKey("kandinsky", &cfg); err != nil {
		slog.Error("kandinsky parse error", slog.Any("err", err))
		panic("")
	}

	slog.Debug("kandinsky config", slog.Any("value", cfg))

	return cfg
}

// Config ...
func (c *KandinskyConfig) Config() *kandinsky.ClientConfig {
	return &kandinsky.ClientConfig{
		BaseURL: c.BaseURL,
		Key:     c.Key,
		Secret:  c.Secret,
	}
}
