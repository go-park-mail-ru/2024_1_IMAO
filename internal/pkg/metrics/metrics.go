package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	totalHits *prometheus.CounterVec
	path      string
	codes     *prometheus.CounterVec
	duration  *prometheus.HistogramVec
}

func CreateMetrics(service string) (*Metrics, error) {
	//var metric Metrics

	//metric.totalHits = prometheus.NewCounterVec()
}
