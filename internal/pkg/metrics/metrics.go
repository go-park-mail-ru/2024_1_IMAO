package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	service  string
	codesNum *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

func CreateMetrics(service string) *Metrics {

}
