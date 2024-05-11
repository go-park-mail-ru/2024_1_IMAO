package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type GRPCMetrics struct {
	totalHits   *prometheus.CounterVec
	serviceName string
	duration    *prometheus.HistogramVec
}

func CreateGRPCMetrics(service string) (*GRPCMetrics, error) {
	var metric GRPCMetrics
	metric.serviceName = service

	metric.totalHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: service + "_total_hits_count",
			Help: "Number of total http requests",
		},
		[]string{"service", "method", "code"})
	if err := prometheus.Register(metric.totalHits); err != nil {
		return nil, err
	}

	metric.duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: service + "_code",
			Help: "Request time",
		},
		[]string{"service", "method", "code"})
	if err := prometheus.Register(metric.duration); err != nil {
		return nil, err
	}

	return &metric, nil
}

func (m *GRPCMetrics) IncreaseTotal(code, method string) {
	m.totalHits.WithLabelValues(m.serviceName, method, code).Inc()
}

func (m *GRPCMetrics) AddDuration(code, method string, duration time.Duration) {
	m.duration.WithLabelValues(m.serviceName, method, code).Observe(duration.Seconds())
}
