package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/WildEgor/pi-storyteller/internal/configs"
)

// PromMetricsRegistry ...
type PromMetricsRegistry struct {
	Reg *prometheus.Registry
}

// NewPromMetricsRegistry ...
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
