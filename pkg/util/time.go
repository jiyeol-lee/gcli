package util

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

// EndOfDayTime function returns the RFC3339 formatted time for the end of the day
func EndOfDayTime() string {
	now := time.Now()

	midnight := time.Date(
		now.Year(), now.Month(), now.Day(),
		23, 59, 59, 0, now.Location(),
	)

	return midnight.Format(time.RFC3339)
}

// ParseUntilStringToTime function parses a string to a time.Time
func ParseUntilStringToTime(until string) (time.Time, error) {
  layout:="20060102T150405Z"
  untilTime, err := time.Parse(layout, until)
  if err != nil {
    return time.Time{}, fmt.Errorf("error parsing until time: %w", err)
  }

  return untilTime, nil
}

// CalculateTimeGap function calculates the duration between two RFC3339 formatted times
func CalculateTimeGap(startTimeRFC3339, endTimeRFC3339 string) (time.Duration, error) {
	startTime, err := time.Parse(time.RFC3339, startTimeRFC3339)
	if err != nil {
		return 0, fmt.Errorf("error parsing start time: %w", err)
	}

	endTime, err := time.Parse(time.RFC3339, endTimeRFC3339)
	if err != nil {
		return 0, fmt.Errorf("error parsing end time: %w", err)
	}

	startTimeOnly := time.Date(
		2000,
		1,
		1,
		startTime.Hour(),
		startTime.Minute(),
		startTime.Second(),
		0,
		time.UTC,
	)
	endTimeOnly := time.Date(
		2000,
		1,
		1,
		endTime.Hour(),
		endTime.Minute(),
		endTime.Second(),
		0,
		time.UTC,
	)

	duration := endTimeOnly.Sub(startTimeOnly)

	return duration, nil
}
