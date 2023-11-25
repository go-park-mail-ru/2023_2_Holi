package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	workTime *prometheus.HistogramVec
	hits     *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	m := &Metrics{
		workTime: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "duration",
			Help:    "Work time for request with certain status.",
			Buckets: prometheus.LinearBuckets(0, 50, 6),
		}, []string{"status", "path"}),
		hits: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "hits",
			Help: "Number of hits.",
		}, []string{"status", "path"}),
	}
	prometheus.MustRegister(m.workTime, m.hits)
	return m
}
