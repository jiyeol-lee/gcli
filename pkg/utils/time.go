package utils

import (
	"fmt"
	"time"
)

// StartOfDayTime function returns the RFC3339 formatted time for the start of the day
func StartOfDayTime() string {
	now := time.Now()

	midnight := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, now.Location(),
	)

	return midnight.Format(time.RFC3339)
}

// CalculateDuration function calculates the duration between two RFC3339 formatted times
func CalculateDuration(startTimeRFC3339, endTimeRFC3339 string) (time.Duration, error) {
	startTime, err := time.Parse(time.RFC3339, startTimeRFC3339)
	if err != nil {
		return 0, fmt.Errorf("error parsing start time: %w", err)
	}

	endTime, err := time.Parse(time.RFC3339, endTimeRFC3339)
	if err != nil {
		return 0, fmt.Errorf("error parsing end time: %w", err)
	}

	duration := endTime.Sub(startTime)

	return duration, nil
}
