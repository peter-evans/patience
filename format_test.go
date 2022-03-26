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

func TestUnifiedDiffTextWithOptions(t *testing.T) {
	type args struct {
		diffs []DiffLine
		opts  UnifiedDiffOptions
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test multiple hunks with context",
			args: args{
				diffs: []DiffLine{
					{Type: Equal, Text: "a"},
					{Type: Equal, Text: "b"},
					{Type: Insert, Text: "c"},
					{Type: Equal, Text: "d"},
					{Type: Equal, Text: "e"},
					{Type: Equal, Text: "f"},
					{Type: Delete, Text: "g"},
					{Type: Insert, Text: "h"},
					{Type: Equal, Text: "i"},
					{Type: Insert, Text: "j"},
					{Type: Equal, Text: "k"},
					{Type: Equal, Text: "l"},
				},
				opts: UnifiedDiffOptions{
					Precontext:  1,
					Postcontext: 1,
				},
			},
			want: "@@ -2,2 +2,3 @@\n b\n+c\n d\n@@ -5,4 +6,5 @@\n f\n-g\n+h\n i\n+j\n k",
		},
		{
			name: "Test source and destination file headers",
			args: args{
				diffs: []DiffLine{
					{Type: Equal, Text: "a"},
					{Type: Equal, Text: "b"},
					{Type: Insert, Text: "c"},
					{Type: Equal, Text: ""},
				},
				opts: UnifiedDiffOptions{
					Precontext:  1,
					Postcontext: 1,
					SrcHeader:   "a.txt",
					DstHeader:   "b.txt",
				},
			},
			want: "--- a.txt\n+++ b.txt\n@@ -2,2 +2,3 @@\n b\n+c\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnifiedDiffTextWithOptions(tt.args.diffs, tt.args.opts); got != tt.want {
				t.Errorf("UnifiedDiffTextWithOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
