package main

import (
	"testing"
	"time"
)

// Função de exemplo que retorna um resultado e um erro
func exampleFunc() (interface{}, error) {
	time.Sleep(100 * time.Millisecond)
	return "Result", nil
	// return nil, errors.New("Error")
}

func TestAdaptiveThrottling(t *testing.T) {
	opts := AdaptiveThrottlingOptions{
		HistoryTimeMinute:    2,
		K:                    2,
		UpperLimitToReject:   0.9,
		MaxRequestDurationMs: 300,
	}

	adaptiveThrottling := AdaptiveThrottling(opts)

	for i := 0; i < 10; i++ {
		result, err := adaptiveThrottling(exampleFunc)
		if err != nil {
			if _, ok := err.(ThrottledException); ok {
				t.Log("Request throttled")
			} else {
				t.Error("Error:", err)
			}
		} else {
			t.Log("Result:", result)
		}
	}
}
