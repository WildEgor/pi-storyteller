package configs

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

var (
	// LogJsonFormat specify json format
	LogJsonFormat string = "json"
)

var logLevelToSlogLevel = map[string]slog.Leveler{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
}

// LoggerConfig holds logger configurations
type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// NewLoggerConfig create logger config
func NewLoggerConfig(c *Configurator) *LoggerConfig {
	cfg := &LoggerConfig{}

	if err := viper.UnmarshalKey("logger", cfg); err != nil {
		slog.Error("logger parse error", slog.Any("err", err))
		panic("logger parse error")
	}

	slog.Info("config", slog.Any("value", cfg))

	so := &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevelToSlogLevel[cfg.Level],
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, so))
	if cfg.Format == LogJsonFormat {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, so))
	}
	slog.SetDefault(logger)

	return cfg
}
