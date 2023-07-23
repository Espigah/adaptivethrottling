package testers

import (
	"fmt"

	"github.com/Espigah/adaptivethrottling"
)

func NewTest1() func() {
	testName := "adaptive_throttling_low_intermittence"

	opts := adaptivethrottling.Options{
		HistoryTimeMinute:    2,
		K:                    2,
		UpperLimitToReject:   0.9,
		MaxRequestDurationMs: 1000,
	}
	throttling := adaptivethrottling.New(opts)

	createGetCommand := func(testName string) func() (interface{}, error) {
		return func() (interface{}, error) {
			return get(testName, &Config{FirstPointOfFailure: 2000, Intermittency: 10000})
		}
	}

	return func() {
		_, err := throttling(createGetCommand(testName))
		if err != nil {
			if _, ok := err.(adaptivethrottling.ThrottledException); ok {
				fmt.Printf("%+v\n", "Request degraded")
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

	opts := adaptivethrottling.Options{
		HistoryTimeMinute:    2,
		K:                    2,
		UpperLimitToReject:   0.9,
		MaxRequestDurationMs: 1000,
	}

	adaptiveThrottling := adaptivethrottling.New(opts)

	createGetCommand := func(testName string) func() (interface{}, error) {
		return func() (interface{}, error) {
			return get(testName, &Config{FirstPointOfFailure: 500, Intermittency: 100})
		}
	}

	return func() {
		_, err := adaptiveThrottling(createGetCommand(testName))
		if err != nil {
			if _, ok := err.(adaptivethrottling.ThrottledException); ok {
				fmt.Printf("%+v\n", "Request degraded")
				m.processed.WithLabelValues(testName, "degraded").Inc()
			} else {
				m.processed.WithLabelValues(testName, "error").Inc()
			}
		} else {
			m.processed.WithLabelValues(testName, "success").Inc()
		}
	}
}
