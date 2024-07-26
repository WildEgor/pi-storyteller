package configs

import (
	"log/slog"

	"github.com/spf13/viper"
)

// MetricsConfig holds metrics config
type MetricsConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// NewMetricsConfig creates metrics config
func NewMetricsConfig() *MetricsConfig {
	cfg := &MetricsConfig{}

	if err := viper.UnmarshalKey("metrics", &cfg); err != nil {
		slog.Error("metrics parse error", slog.Any("err", err))
		panic("")
	}

	slog.Debug("metrics config", slog.Any("value", cfg))

	return cfg
}
