package testers

import (
	"github.com/sony/gobreaker"
)

func NewTest3() func() {
	testName := "circuit_breaker_low_intermittence"

	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	cb := gobreaker.NewCircuitBreaker(st)

	createGetCommand := func(testName string) func() (interface{}, error) {
		return func() (interface{}, error) {
			return get(testName, &Config{FirstPointOfFailure: 2000, Intermittency: 2000})
		}
	}

	return func() {
		_, err := cb.Execute(createGetCommand(testName))
		if err != nil {

			if err.Error() == "circuit breaker is open" {
				m.processed.WithLabelValues(testName, "degraded").Inc()
			} else {
				m.processed.WithLabelValues(testName, "error").Inc()
			}
		} else {
			m.processed.WithLabelValues(testName, "success").Inc()
		}
	}
}

func NewTest4() func() {
	testName := "circuit_breaker_high_intermittence"

	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	cb := gobreaker.NewCircuitBreaker(st)

	createGetCommand := func(testName string) func() (interface{}, error) {
		return func() (interface{}, error) {
			return get(testName, &Config{FirstPointOfFailure: 2000, Intermittency: 10})
		}
	}

	return func() {
		_, err := cb.Execute(createGetCommand(testName))
		if err != nil {

			if err.Error() == "circuit breaker is open" {
				m.processed.WithLabelValues(testName, "degraded").Inc()
			} else {
				m.processed.WithLabelValues(testName, "error").Inc()
			}
		} else {
			m.processed.WithLabelValues(testName, "success").Inc()
		}
	}
}
