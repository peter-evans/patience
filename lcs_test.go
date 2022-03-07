// Package patience implements the Patience Diff algorithm.
package patience

import (
	"reflect"
	"testing"
)

func TestLCS(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want [][2]int
	}{
		{
			name: "Test identical slices",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"a", "b", "c"},
			},
			want: [][2]int{
				{0, 0},
				{1, 1},
				{2, 2},
			},
		},
		{
			name: "Test slices of different lengths",
			args: args{
				a: []string{"a", "z", "b", "c"},
				b: []string{"a", "b", "y", "w", "c"},
			},
			want: [][2]int{
				{0, 0},
				{2, 1},
				{3, 4},
			},
		},
		{
			name: "Test slices containing duplicate elements",
			args: args{
				a: []string{"a", "b", "a", "y", "c", "c"},
				b: []string{"z", "b", "a", "c", "c", "b"},
			},
			want: [][2]int{
				{1, 1},
				{2, 2},
				{4, 3},
				{5, 4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LCS(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LCS() = %v, want %v", got, tt.want)
			}
		})
	}
}
