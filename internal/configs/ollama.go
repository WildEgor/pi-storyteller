package configs

import (
	"log/slog"

	"github.com/spf13/viper"
)

// KandinskyConfig holds kandinsky api configuration
type OllamaConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

// NewOllamaConfig creates kandinsky config
func NewOllamaConfig() *OllamaConfig {
	cfg := &OllamaConfig{}

	if err := viper.UnmarshalKey("ollama", &cfg); err != nil {
		panic("ollama config parse error")
	}

	slog.Info("config", slog.Any("value", cfg))

	return cfg
}
