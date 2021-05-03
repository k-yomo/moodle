package moodle

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestWarnings_Error(t *testing.T) {
	tests := []struct {
		name string
		l    Warnings
		want string
	}{
		{
			l: []*Warning{
				{
					Item:        "quiz",
					ItemID:      1111,
					WarningCode: "1",
					Message:     "This quiz is not currently available",
				},
			},
			want: "item: quiz, itemID: 1111, warningCode: 1, message: This quiz is not currently available",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(tt.l.Error(), tt.want); diff != "" {
				t.Errorf("Error() (-got, +want)\n%s", diff)
			}
		})
	}
}
