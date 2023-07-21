package testers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	processed *prometheus.CounterVec
}

var (
	m = newMetrics()
)

func newMetrics() *metrics {

	m := &metrics{
		processed: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "operations_count",
			Help: "Number of requests.",
		}, []string{"name", "status"}),
	}

	return m
}
