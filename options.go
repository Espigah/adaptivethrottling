package main

import (
	"os"
	"strconv"
)

// Default values
const (
	K                       = "2"   // Multiplier that determines the aggressiveness of throttling
	HISTORY_TIME_MINUTE     = "120" // Time in minutes to keep request history
	UPPER_LIMIT_TO_REJECT   = "0.9" // Upper limit to reject new requests
	MAX_REQUEST_DURATION_MS = "300" // Maximum request duration in milliseconds
)

// AdaptiveThrottlingOptions contém as opções para o adaptive throttling.
type AdaptiveThrottlingOptions struct {
	HistoryTimeMinute    int
	K                    float64
	UpperLimitToReject   float64
	MaxRequestDurationMs int
}

func (o *AdaptiveThrottlingOptions) Fill() {
	if o.HistoryTimeMinute == 0 {
		o.HistoryTimeMinute, _ = strconv.Atoi(o.GetOrDefault("HISTORY_TIME_MINUTE", HISTORY_TIME_MINUTE))
	}

	if o.K == 0 {
		o.K, _ = strconv.ParseFloat(o.GetOrDefault("K", K), 64)
	}

	if o.UpperLimitToReject == 0 {
		o.UpperLimitToReject, _ = strconv.ParseFloat(o.GetOrDefault("UPPER_LIMIT_TO_REJECT", UPPER_LIMIT_TO_REJECT), 64)
	}

	if o.MaxRequestDurationMs == 0 {
		o.MaxRequestDurationMs, _ = strconv.Atoi(o.GetOrDefault("MAX_REQUEST_DURATION_MS", MAX_REQUEST_DURATION_MS))
	}
}

func (o *AdaptiveThrottlingOptions) GetOrDefault(envVar, defaultValue string) string {
	value, exists := os.LookupEnv(envVar)
	if exists {
		return value
	}
	return defaultValue
}
