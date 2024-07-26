package configs

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Configurator ...
type Configurator struct {
	watchers []func()
}

// NewConfigurator create new Configurator
func NewConfigurator() *Configurator {
	c := &Configurator{
		watchers: make([]func(), 0),
	}
	c.load()
	return c
}

// Watch config
func (c *Configurator) Watch() {
	viper.OnConfigChange(func(e fsnotify.Event) {
		slog.Info(fmt.Sprintf("watchers len: %d", len(c.watchers)))

		for _, watcher := range c.watchers {
			watcher()
		}
	})
	viper.WatchConfig()
}

// Register watcher
func (c *Configurator) Register(name string, fn func()) {
	slog.Info("register watcher", slog.Any("value", name))
	c.watchers = append(c.watchers, fn)
}

// load Load env data from files (default: .env, .env.local)
func (c *Configurator) load() {
	path := os.Getenv("PI_STORYTELLER_CONFIG_PATH")
	if len(path) != 0 {
		viper.AddConfigPath(path)
	} else {
		//nolint
		pwd, _ := os.Getwd()
		viper.AddConfigPath(pwd)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("error loading config file", slog.Any("err", err))
		panic("error loading config file")
	}

	if err := viper.MergeInConfig(); err != nil {
		slog.Error("error merge config file", slog.Any("err", err))
		return
	}
}
