package util

import (
	"fmt"
	"reflect"
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

func TestEndOfDayTime(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Is RFC3339 formatted time for the end of the day",
			want: func() string {
				now := time.Now()

				return fmt.Sprintf(
					"^%s-%02d-%sT%s:%s:%s-\\d{2}:\\d{2}$",
					strconv.Itoa(now.Year()),
					now.Month(),
					strconv.Itoa(now.Day()),
					"23",
					"59",
					"59",
				)
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EndOfDayTime()
			match, matchErr := regexp.MatchString(tt.want, got)
			if matchErr != nil {
				t.Errorf("EndOfDayTime() error = %v", matchErr)
				return
			}
			if !match {
				t.Errorf("EndOfDayTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestParseUntilStringToRFC3339(t *testing.T) {
// 	type args struct {
// 		until string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
//     {
//       name: "Correct date format (1)",
//       args: args{
//         until: "20231002T035959Z",
//       },
//       want: "2023-10-02T03:59:59Z",
//       wantErr: false,
//     },
//     {
//       name: "Correct date format (2)",
//       args: args{
//         until: "20240115T123000Z",
//       },
//       want: "2024-01-15T12:30:00Z",
//       wantErr: false,
//     },
//     {
//       name: "Correct date format (3)",
//       args: args{
//         until: "20250720T084500Z",
//       },
//       want: "2025-07-20T08:45:00Z",
//       wantErr: false,
//     },
//     {
//       name: "Wrong date format (1)",
//       args: args{
//         until: "2020202020202020",
//       },
//       want: "",
//       wantErr: true,
//     },
//     {
//       name: "Wrong date format (2)",
//       args: args{
//         until: "",
//       },
//       want: "",
//       wantErr: true,
//     },
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := ParseUntilStringToTime(tt.args.until)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ParseUntilStringToTime() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("ParseUntilStringToTime() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestParseUntilStringToTime(t *testing.T) {
	type args struct {
		until string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
    {
      name: "Correct date format (1)",
      args: args{
        until: "20231002T035959Z",
      },
      want: time.Date(2023, 10, 2, 3, 59, 59, 0, time.UTC),
      wantErr: false,
    },
    {
      name: "Correct date format (2)",
      args: args{
        until: "20240115T123000Z",
      },
      want: time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC),
      wantErr: false,
    },
    {
      name: "Correct date format (3)",
      args: args{
        until: "20250720T084500Z",
      },
      want: time.Date(2025, 7, 20, 8, 45, 0, 0, time.UTC),
      wantErr: false,
    },
    {
      name: "Wrong date format (1)",
      args: args{
        until: "2020202020202020",
      },
      want: time.Time{},
      wantErr: true,
    },
    {
      name: "Wrong date format (2)",
      args: args{
        until: "",
      },
      want: time.Time{},
      wantErr: true,
    },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUntilStringToTime(tt.args.until)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUntilStringToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseUntilStringToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateTimeGap(t *testing.T) {
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
		{
			name: "Only time is considered, date is ignored",
			args: args{
				startTimeRFC3339: "2021-01-01T00:00:00Z",
				endTimeRFC3339:   "2031-12-31T02:02:02Z",
			},
			want:    time.Hour*2 + time.Minute*2 + time.Second*2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateTimeGap(tt.args.startTimeRFC3339, tt.args.endTimeRFC3339)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateTimeGap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateTimeGap() = %v, want %v", got, tt.want)
			}
		})
	}
}


