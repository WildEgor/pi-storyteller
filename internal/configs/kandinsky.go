package configs

import (
	"log/slog"

	"github.com/WildEgor/pi-storyteller/pkg/kandinsky"
	"github.com/spf13/viper"
)

// KandinskyConfig holds kandinsky api configuration
type KandinskyConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Key     string `mapstructure:"key"`
	Secret  string `mapstructure:"secret"`
}

// NewKandinskyConfig creates kandinsky config
func NewKandinskyConfig(c *Configurator) *KandinskyConfig {
	cfg := &KandinskyConfig{}

	if err := viper.UnmarshalKey("kandinsky", &cfg); err != nil {
		panic("kandinsky config parse error")
	}

	slog.Info("config", slog.Any("value", cfg))

	return cfg
}

func (c *KandinskyConfig) Config() *kandinsky.ClientConfig {
	return &kandinsky.ClientConfig{
		BaseURL: c.BaseURL,
		Key:     c.Key,
		Secret:  c.Secret,
	}
}
