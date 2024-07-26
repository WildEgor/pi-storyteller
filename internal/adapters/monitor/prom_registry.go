package monitor

import (
	"github.com/WildEgor/pi-storyteller/internal/configs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type PromMetricsRegistry struct {
	Reg *prometheus.Registry
}

func NewPromMetricsRegistry(config *configs.MetricsConfig) *PromMetricsRegistry {
	if config.Enabled {
		reg := prometheus.NewRegistry()
		reg.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)
		return &PromMetricsRegistry{
			Reg: reg,
		}
	}

	return &PromMetricsRegistry{}
}
