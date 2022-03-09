// Package patience implements the Patience Diff algorithm.
package patience

import (
	"testing"
)

func TestDiffText(t *testing.T) {
	type args struct {
		diffs []DiffLine
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestDiffText",
			args: args{
				diffs: []DiffLine{
					{
						Type: Equal,
						Text: "a",
					},
					{
						Type: Insert,
						Text: "b",
					},
					{
						Type: Equal,
						Text: "c",
					},
					{
						Type: Equal,
						Text: "",
					},
					{
						Type: Delete,
						Text: "d",
					},
					{
						Type: Equal,
						Text: "e",
					},
				},
			},
			want: " a\n+b\n c\n\n-d\n e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffText(tt.args.diffs); got != tt.want {
				t.Errorf("DiffText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiffTextA(t *testing.T) {
	type args struct {
		diffs []DiffLine
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestDiffText",
			args: args{
				diffs: []DiffLine{
					{
						Type: Equal,
						Text: "a",
					},
					{
						Type: Insert,
						Text: "b",
					},
					{
						Type: Equal,
						Text: "c",
					},
					{
						Type: Equal,
						Text: "",
					},
					{
						Type: Delete,
						Text: "d",
					},
					{
						Type: Equal,
						Text: "e",
					},
				},
			},
			want: " a\n c\n\n-d\n e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffTextA(tt.args.diffs); got != tt.want {
				t.Errorf("DiffTextA() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiffTextB(t *testing.T) {
	type args struct {
		diffs []DiffLine
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestDiffText",
			args: args{
				diffs: []DiffLine{
					{
						Type: Equal,
						Text: "a",
					},
					{
						Type: Insert,
						Text: "b",
					},
					{
						Type: Equal,
						Text: "c",
					},
					{
						Type: Equal,
						Text: "",
					},
					{
						Type: Equal,
						Text: "e",
					},
				},
			},
			want: " a\n+b\n c\n\n e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DiffTextB(tt.args.diffs); got != tt.want {
				t.Errorf("DiffTextB() = %v, want %v", got, tt.want)
			}
		})
	}
}
