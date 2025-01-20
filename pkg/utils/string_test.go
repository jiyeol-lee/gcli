package util

import "testing"

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
