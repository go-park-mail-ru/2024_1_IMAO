package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type Metrics struct {
	totalHits   *prometheus.CounterVec
	serviceName string
	duration    *prometheus.HistogramVec
}

func CreateMetrics(service string) (*Metrics, error) {
	var metric Metrics
	metric.serviceName = service

	metric.totalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: service + "_total_hits_count",
			Help: "Number of total http requests",
		},
		[]string{"path", "service", "code"})
	if err := prometheus.Register(metric.totalHits); err != nil {
		return nil, err
	}

	metric.duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: service + "_code_count",
			Help: "Request time",
		},
		[]string{"path", "service", "code"})
	if err := prometheus.Register(metric.duration); err != nil {
		return nil, err
	}

	return &metric, nil
}

func (m *Metrics) IncreaseTotal(path, code string) {
	m.totalHits.WithLabelValues(m.serviceName, path, code).Inc()
}

func (m *Metrics) AddDuration(path, code string, duration time.Duration) {
	m.duration.WithLabelValues(m.serviceName, path, code).Observe(duration.Seconds())
}
