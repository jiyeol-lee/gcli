package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func TestStartOfDayTime(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Is RFC3339 formatted time for the start of the day",
			want: func() string {
				now := time.Now()

				return fmt.Sprintf(
					"^%s-%02d-%sT%s:%s:%s-\\d{2}:\\d{2}$",
					strconv.Itoa(now.Year()),
					now.Month(),
					strconv.Itoa(now.Day()),
					"00",
					"00",
					"00",
				)
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StartOfDayTime()
			match, matchErr := regexp.MatchString(tt.want, got)
			if matchErr != nil {
				t.Errorf("StartOfDayTime() error = %v", matchErr)
				return
			}
			if !match {
				t.Errorf("StartOfDayTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateDuration(t *testing.T) {
	type args struct {
		startTimeRFC3339 string
		endTimeRFC3339   string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name: "When start time is before end time, duration is positive",
			args: args{
				startTimeRFC3339: "2021-01-01T00:00:00Z",
				endTimeRFC3339:   "2021-01-01T02:00:00Z",
			},
			want:    time.Hour * 2,
			wantErr: false,
		},
		{
			name: "When start time is after end time, duration is negative",
			args: args{
				startTimeRFC3339: "2021-01-01T02:00:00Z",
				endTimeRFC3339:   "2021-01-01T00:00:00Z",
			},
			want:    time.Hour * -2,
			wantErr: false,
		},
		{
			name: "When start time and end time are the same, duration is 0",
			args: args{
				startTimeRFC3339: "2021-01-01T00:00:00Z",
				endTimeRFC3339:   "2021-01-01T00:00:00Z",
			},
			want:    time.Hour * 0,
			wantErr: false,
		},
		{
			name: "When start time is invalid, return error",
			args: args{
				startTimeRFC3339: "invalid",
				endTimeRFC3339:   "2021-01-01T00:00:00Z",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "When end time is invalid, return error",
			args: args{
				startTimeRFC3339: "2021-01-01T00:00:00Z",
				endTimeRFC3339:   "invalid",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateDuration(tt.args.startTimeRFC3339, tt.args.endTimeRFC3339)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
