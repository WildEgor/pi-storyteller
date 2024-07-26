package configs

import (
	"github.com/spf13/viper"

	"log/slog"

	"github.com/WildEgor/pi-storyteller/pkg/chatgpt"
)

// ChatGPTConfig holds ChatGPT api configuration
type ChatGPTConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Key     string `mapstructure:"key"`
}

// NewChatGPTConfig creates ChatGPT config
func NewChatGPTConfig() *ChatGPTConfig {
	cfg := &ChatGPTConfig{}

	if err := viper.UnmarshalKey("chatgpt", &cfg); err != nil {
		slog.Error("chatgpt parse error", slog.Any("err", err))
		panic("")
	}

	slog.Debug("chatgpt config", slog.Any("value", cfg))

	return cfg
}

// Config ...
func (c *ChatGPTConfig) Config() *chatgpt.ClientConfig {
	return &chatgpt.ClientConfig{
		BaseURL: c.BaseURL,
		Key:     c.Key,
	}
}
