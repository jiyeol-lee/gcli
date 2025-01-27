package util

import (
	"testing"
)

func TestTruncateWithSuffix(t *testing.T) {
	type args struct {
		s         string
		maxLength int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "When string is shorter than max length, return string",
			args: args{
				s:         "short",
				maxLength: 10,
			},
			want: "short",
		},
		{
			name: "When string is longer than max length, return truncated string with suffix",
			args: args{
				s:         "long string",
				maxLength: 7,
			},
			want: "long st...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateWithSuffix(tt.args.s, tt.args.maxLength); got != tt.want {
				t.Errorf("TruncateWithSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findUntilFromRecurrence(t *testing.T) {
	type args struct {
		recurrence []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
    {
      name: "When recurrence is empty, return empty string",
      args: args{
        recurrence: []string{},
      },
      want: "",
    },
    {
      name: "When recurrence is nil, return empty string",
      args: args{
        recurrence: nil,
      },
      want: "",
    },
    {
      name: "When recurrence has UNTIL, return UNTIL value (1)",
      args: args{
        recurrence: []string{
          "RRULE:FREQ=WEEKLY;WKST=SU;UNTIL=20240115T045959Z;INTERVAL=2;BYDAY=MO",
        },
      },
      want: "20240115T045959Z",
    },
    {
      name: "When recurrence has UNTIL, return UNTIL value (2)",
      args: args{
        recurrence: []string{
          "RRULE:FREQ=WEEKLY;WKST=SU;UNTIL=20231002T035959Z;BYDAY=MO,TU,WE",
        },
      },
      want: "20231002T035959Z",
    },
    {
      name: "When recurrence has no UNTIL, return empty string (1)",
      args: args{
        recurrence: []string{
          "RRULE:FREQ=WEEKLY;WKST=SU;INTERVAL=2;BYDAY=MO",
        },
      },
      want: "",
    },
    {
      name: "When recurrence has no UNTIL, return empty string (2)",
      args: args{
        recurrence: []string{
          "RRULE:FREQ=WEEKLY;WKST=SU;BYDAY=MO,TU,WE",
        },
      },
      want: "",
    },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindUntilFromRecurrence(tt.args.recurrence); got != tt.want {
				t.Errorf("findUntilFromRecurrence() = %v, want %v", got, tt.want)
			}
		})
	}
}
