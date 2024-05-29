package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type DatabaseMetrics struct {
	errorsCount *prometheus.CounterVec
	serviceName string
	dbName      string
	duration    *prometheus.HistogramVec
}

func CreateDatabaseMetrics(service, db string) (*DatabaseMetrics, error) {
	var metric DatabaseMetrics
	metric.serviceName = service
	metric.dbName = db

	metric.errorsCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%s_%s_errors_count", db, service),
			Help: "Number of database errors in service",
		},
		[]string{"function", "service", "database"})
	if err := prometheus.Register(metric.errorsCount); err != nil {
		return nil, err
	}

	metric.duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: fmt.Sprintf("%s_%s_total_requests", db, service),
			Help: "Total database requests in service",
		},
		[]string{"function", "service", "database"})
	if err := prometheus.Register(metric.duration); err != nil {
		return nil, err
	}

	return &metric, nil
}

func (m *DatabaseMetrics) IncreaseErrors(funcName string) {
	m.errorsCount.WithLabelValues(funcName, m.serviceName, m.dbName).Inc()
}

func (m *DatabaseMetrics) AddDuration(funcName string, duration time.Duration) {
	m.duration.WithLabelValues(funcName, m.serviceName, m.dbName).Observe(duration.Seconds())
}
