package testers

import (
	adaptivethrottlinggo "github.com/Espigah/adaptive-throttling-go"
)

func NewTest1() func() {
	testName := "adaptive_throttling_low_intermittence"

	opts := adaptivethrottlinggo.Options{
		HistoryTimeMinute:    2,
		K:                    2,
		UpperLimitToReject:   0.9,
		MaxRequestDurationMs: 300,
	}
	adaptiveThrottling := adaptivethrottlinggo.New(opts)

	createGetCommand := func(testName string) func() (interface{}, error) {
		return func() (interface{}, error) {
			return get(testName, &Config{FirstPointOfFailure: 2000, Intermittency: 2000})
		}
	}

	return func() {
		_, err := adaptiveThrottling(createGetCommand(testName))
		if err != nil {
			if _, ok := err.(adaptivethrottlinggo.ThrottledException); ok {
				//t.Log("Request throttled")
				m.processed.WithLabelValues(testName, "degraded").Inc()
			} else {
				m.processed.WithLabelValues(testName, "error").Inc()
			}
		} else {
			m.processed.WithLabelValues(testName, "success").Inc()
		}
	}
}

func NewTest2() func() {
	testName := "adaptive_throttling_high_intermittence"

	opts := adaptivethrottlinggo.Options{
		HistoryTimeMinute:    2,
		K:                    2,
		UpperLimitToReject:   0.9,
		MaxRequestDurationMs: 300,
	}

	adaptiveThrottling := adaptivethrottlinggo.New(opts)

	createGetCommand := func(testName string) func() (interface{}, error) {
		return func() (interface{}, error) {
			return get(testName, &Config{FirstPointOfFailure: 2000, Intermittency: 10})
		}
	}

	return func() {
		_, err := adaptiveThrottling(createGetCommand(testName))
		if err != nil {
			if _, ok := err.(adaptivethrottlinggo.ThrottledException); ok {
				//t.Log("Request throttled")
				m.processed.WithLabelValues(testName, "degraded").Inc()
			} else {
				m.processed.WithLabelValues(testName, "error").Inc()
			}
		} else {
			m.processed.WithLabelValues(testName, "success").Inc()
		}
	}
}
