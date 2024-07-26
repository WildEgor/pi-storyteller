package configs

import (
	"github.com/spf13/viper"

	"log/slog"
)

// OllamaConfig holds Ollama api configuration
type OllamaConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

// NewOllamaConfig creates kandinsky config
func NewOllamaConfig() *OllamaConfig {
	cfg := &OllamaConfig{}

	if err := viper.UnmarshalKey("ollama", &cfg); err != nil {
		slog.Error("ollama parse error", slog.Any("err", err))
		panic("")
	}

	slog.Debug("ollama config", slog.Any("value", cfg))

	return cfg
}
