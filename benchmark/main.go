package main

import (
	"io/ioutil"
	"net/http"

	a "github.com/Espigah/adaptive-throttling-go"
)

type metrics struct {
	info          *prometheus.GaugeVec
	req_processed *prometheus.CounterVec
	sql_duration  *prometheus.HistogramVec
	rec_duration  *prometheus.HistogramVec
	api_duration  *prometheus.HistogramVec
}

var (
	reg = prometheus.NewRegistry()
	m   = NewMetrics(reg)
)

func NewMetrics(reg prometheus.Registerer) *metrics {

	m := &metrics{
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "txn_count",
			Help: "Target amount for completed requests",
		}, []string{"batch"}),

		req_processed: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "lb_fs_etl_operations_count",
			Help: "Number of completed requests.",
		}, []string{"batch"}),

		sql_duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "lb_fs_sql_duration_seconds",
			Help: "Duration of the sql requests",
			// 4 times larger apdex status
			// Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
			// Buckets: prometheus.LinearBuckets(0.1, 5, 5),
			Buckets: []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"batch"}),

		rec_duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "lb_fs_etl_operations_seconds",
			Help: "Duration of the entire requests",

			Buckets: []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"batch"}),

		api_duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "lb_fs_api_duration_seconds",
			Help:    "Duration of the api requests",
			Buckets: []float64{0.1, 0.15, 0.2, 0.25, 0.3},
		}, []string{"batch"}),
	}

	reg.MustRegister(m.info, m.req_processed, m.sql_duration, m.rec_duration, m.api_duration)

	return m
}


func main() {
	opts := a.AdaptiveThrottlingOptions{
		HistoryTimeMinute:    2,
		K:                    2,
		UpperLimitToReject:   0.9,
		MaxRequestDurationMs: 300,
	}
}

var cb *breaker.CircuitBreaker

func Get(url string) ([]byte, error) {
	body, err := cb.Execute(HTTPGet(url))
	if err != nil {
		return nil, err
	}

	return body.([]byte), nil
}

func HTTPGet(url string) interface{} {
	return func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	}
}
