package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/WildEgor/pi-storyteller/internal/configs"
)

var _ Monitor = (*PromMetrics)(nil)

const (
	labelApp      = "app"
	labelKind     = "kind"
	labelUsername = "username"
)

const (
	// ProblemKindUnknown ...
	ProblemKindUnknown = "unknown"
	// ProblemKindParse ...
	ProblemKindParse = "parse"
	// ProblemKindDecode ...
	ProblemKindDecode = "decode"
	// ProblemKindPublish ...
	ProblemKindPublish = "send"
)

// PromMetrics ...
type PromMetrics struct {
	appConfig                      *configs.AppConfig
	metricsConfig                  *configs.MetricsConfig
	completedJobs, problematicJobs *prometheus.CounterVec
	activeJobs                     prometheus.Gauge
	regAdapter                     *PromMetricsRegistry
}

// NewPromMetrics ...
func NewPromMetrics(reg *PromMetricsRegistry, appConfig *configs.AppConfig, metricsConfig *configs.MetricsConfig) *PromMetrics {
	return &PromMetrics{
		appConfig:     appConfig,
		metricsConfig: metricsConfig,
		regAdapter:    reg,
		completedJobs: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "completed_jobs_total",
			Help: "The total number of completed jobs",
		},
			[]string{labelApp, labelUsername},
		),
		problematicJobs: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "problematic_jobs_total",
			Help: "The total number of failed jobs",
		},
			[]string{labelApp, labelUsername, labelKind},
		),
		activeJobs: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "active_jobs_current",
			Help: "Current number of active jobs",
		}),
	}
}

// IncActiveJobsCounter ...
func (m *PromMetrics) IncActiveJobsCounter() {
	if m.metricsConfig.Enabled {
		m.activeJobs.Inc()
	}
}

// DecActiveJobsCounter ...
func (m *PromMetrics) DecActiveJobsCounter() {
	if m.metricsConfig.Enabled {
		m.activeJobs.Dec()
	}
}

// IncAllJobsCounter ...
func (m *PromMetrics) IncAllJobsCounter(username string) {
	if m.metricsConfig.Enabled {
		m.completedJobs.With(prometheus.Labels{labelApp: m.appConfig.Name, labelUsername: username}).Inc()
	}
}

// IncFailedJobsCounter ...
func (m *PromMetrics) IncFailedJobsCounter(username, kind string) {
	if m.metricsConfig.Enabled {
		m.problematicJobs.With(prometheus.Labels{labelApp: m.appConfig.Name, labelUsername: username, labelKind: kind}).Inc()
	}
}
