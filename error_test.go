package moodle

import (
	"errors"
	"testing"
)

func TestCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "when APIError, it returns the error code",
			args: args{err: &APIError{ErrorCode: "invalidtoken"}},
			want: "invalidtoken",
		},
		{
			name: "when not APIError, it returns unknown",
			args: args{err: errors.New("invalid")},
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Code(tt.args.err); got != tt.want {
				t.Errorf("Code() = %v, want %v", got, tt.want)
			}
		})
	}
}
