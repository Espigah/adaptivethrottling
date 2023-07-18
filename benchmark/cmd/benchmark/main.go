package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	adaptivethrottlinggo "github.com/Espigah/adaptive-throttling-go"
	"github.com/prometheus/client_golang/prometheus"
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

const api = "http://localhost:3000/"

type Config struct {
	FirstPointOfFailure int `json: "firstPointOfFailure"`
	Intermittency       int `json: "intermittency"`
}

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
	opts := adaptivethrottlinggo.AdaptiveThrottlingOptions{
		HistoryTimeMinute:    2,
		K:                    2,
		UpperLimitToReject:   0.9,
		MaxRequestDurationMs: 300,
	}

	adaptiveThrottling := adaptivethrottlinggo.AdaptiveThrottling(opts)
	for i := 0; i < 10; i++ {
		_, err := adaptiveThrottling(exampleFunc)
		if err != nil {
			if _, ok := err.(adaptivethrottlinggo.ThrottledException); ok {
				//t.Log("Request throttled")
			} else {
				//t.Error("Error:", err)
			}
		} else {
			//t.Log("Result:", result)
		}
	}
}

// Função de exemplo que retorna um resultado e um erro
func exampleFunc() (interface{}, error) {
	Get("teste1", &Config{FirstPointOfFailure: 5, Intermittency: 2})
	return "Result", nil
	// return nil, errors.New("Error")
}

//var cb *breaker.CircuitBreaker

func Get(testName string, body *Config) (interface{}, error) {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		log.Fatal(err)
	}

	url := api + testName

	req, err := http.NewRequest(http.MethodGet, url, &buf)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	return res, err
}
