package patience

import (
	"reflect"
	"testing"
)

func Test_makeHunks(t *testing.T) {
	e := DiffLine{
		Type: Equal,
		Text: "e",
	}
	i := DiffLine{
		Type: Insert,
		Text: "i",
	}
	d := DiffLine{
		Type: Delete,
		Text: "d",
	}

	type args struct {
		diffs       []DiffLine
		precontext  int
		postcontext int
	}
	tests := []struct {
		name string
		args args
		want []Hunk
	}{
		{
			name: "Test nil diffs",
			args: args{
				diffs:       nil,
				precontext:  2,
				postcontext: 2,
			},
			want: nil,
		},
		{
			name: "Test no diff (all equalities)",
			args: args{
				diffs: []DiffLine{
					e, e, e, e,
				},
				precontext:  2,
				postcontext: 2,
			},
			want: nil,
		},
		{
			name: "Test all modifications (no equalities)",
			args: args{
				diffs: []DiffLine{
					i, i, i, i, d, d, d, d,
				},
				precontext:  2,
				postcontext: 2,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{i, i, i, i, d, d, d, d},
					SrcStart: 1,
					SrcLines: 4,
					DstStart: 1,
					DstLines: 4,
				},
			},
		},
		{
			name: "Test deletions only",
			args: args{
				diffs: []DiffLine{
					d, d, d, d,
				},
				precontext:  2,
				postcontext: 2,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{d, d, d, d},
					SrcStart: 1,
					SrcLines: 4,
					DstStart: 0,
					DstLines: 0,
				},
			},
		},
		{
			name: "Test insertions only",
			args: args{
				diffs: []DiffLine{
					i, i, i, i,
				},
				precontext:  2,
				postcontext: 2,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{i, i, i, i},
					SrcStart: 0,
					SrcLines: 0,
					DstStart: 1,
					DstLines: 4,
				},
			},
		},
		{
			name: "Test deletions and equalities only",
			args: args{
				diffs: []DiffLine{
					d, d, e, e, e,
				},
				precontext:  2,
				postcontext: 2,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{d, d, e, e},
					SrcStart: 1,
					SrcLines: 4,
					DstStart: 1,
					DstLines: 2,
				},
			},
		},
		{
			name: "Test insertions and equalities only",
			args: args{
				diffs: []DiffLine{
					i, i, e, e, e,
				},
				precontext:  2,
				postcontext: 2,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{i, i, e, e},
					SrcStart: 1,
					SrcLines: 2,
					DstStart: 1,
					DstLines: 4,
				},
			},
		},
		{
			name: "Test updating DstStart when the current hunk has no destination diff lines",
			args: args{
				diffs: []DiffLine{
					d, d, e, e, e, e, e, e, e, d, d,
				},
				precontext:  3,
				postcontext: 3,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{d, d, e, e, e},
					SrcStart: 1,
					SrcLines: 5,
					DstStart: 1,
					DstLines: 3,
				},
				{
					Diffs:    []DiffLine{e, e, e, d, d},
					SrcStart: 7,
					SrcLines: 5,
					DstStart: 5,
					DstLines: 3,
				},
			},
		},
		{
			name: "Test updating SrcStart when the current hunk has no source diff lines",
			args: args{
				diffs: []DiffLine{
					i, i, e, e, e, e, e, e, e, i, i,
				},
				precontext:  3,
				postcontext: 3,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{i, i, e, e, e},
					SrcStart: 1,
					SrcLines: 3,
					DstStart: 1,
					DstLines: 5,
				},
				{
					Diffs:    []DiffLine{e, e, e, i, i},
					SrcStart: 5,
					SrcLines: 3,
					DstStart: 7,
					DstLines: 5,
				},
			},
		},
		{
			name: "Test equal block sizes greater than context",
			args: args{
				diffs: []DiffLine{
					e, e, e, d, i, i, e, e, e, e, e, d, i, e, e, e,
				},
				precontext:  2,
				postcontext: 2,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{e, e, d, i, i, e, e},
					SrcStart: 2,
					SrcLines: 5,
					DstStart: 2,
					DstLines: 6,
				},
				{
					Diffs:    []DiffLine{e, e, d, i, e, e},
					SrcStart: 8,
					SrcLines: 5,
					DstStart: 9,
					DstLines: 5,
				},
			},
		},
		{
			name: "Test equal block sizes less than context",
			args: args{
				diffs: []DiffLine{
					e, d, i, e, d, i, e,
				},
				precontext:  2,
				postcontext: 2,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{e, d, i, e, d, i, e},
					SrcStart: 1,
					SrcLines: 5,
					DstStart: 1,
					DstLines: 5,
				},
			},
		},
		{
			name: "Test the maximum equal block size within a hunk",
			args: args{
				diffs: []DiffLine{
					d, i, e, e, e, e, d, i,
				},
				precontext:  2,
				postcontext: 2,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{d, i, e, e, e, e, d, i},
					SrcStart: 1,
					SrcLines: 6,
					DstStart: 1,
					DstLines: 6,
				},
			},
		},
		{
			name: "Test multiple modified blocks within a hunk",
			args: args{
				diffs: []DiffLine{
					d, i, e, e, e, e, d, d, i, e, e, e, d, i, e, e, e, e, e, d, i, e, e, e,
				},
				precontext:  2,
				postcontext: 2,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{d, i, e, e, e, e, d, d, i, e, e, e, d, i, e, e},
					SrcStart: 1,
					SrcLines: 13,
					DstStart: 1,
					DstLines: 12,
				},
				{
					Diffs:    []DiffLine{e, e, d, i, e, e},
					SrcStart: 15,
					SrcLines: 5,
					DstStart: 14,
					DstLines: 5,
				},
			},
		},
		{
			name: "Test differing pre and post contexts",
			args: args{
				diffs: []DiffLine{
					d, i, e, e, e, e, d, d, i, e, e, e, d, d, i, e, e, e, e, e, d, i, e, e, e,
				},
				precontext:  2,
				postcontext: 1,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{d, i, e},
					SrcStart: 1,
					SrcLines: 2,
					DstStart: 1,
					DstLines: 2,
				},
				{
					Diffs:    []DiffLine{e, e, d, d, i, e, e, e, d, d, i, e},
					SrcStart: 4,
					SrcLines: 10,
					DstStart: 4,
					DstLines: 8,
				},
				{
					Diffs:    []DiffLine{e, e, d, i, e},
					SrcStart: 16,
					SrcLines: 4,
					DstStart: 14,
					DstLines: 4,
				},
			},
		},
		{
			name: "Test no context",
			args: args{
				diffs: []DiffLine{
					d, i, e, e, e, e, d, i, e, e, e,
				},
				precontext:  0,
				postcontext: 0,
			},
			want: []Hunk{
				{
					Diffs:    []DiffLine{d, i},
					SrcStart: 1,
					SrcLines: 1,
					DstStart: 1,
					DstLines: 1,
				},
				{
					Diffs:    []DiffLine{d, i},
					SrcStart: 6,
					SrcLines: 1,
					DstStart: 6,
					DstLines: 1,
				},
			},
		},
		{
			name: "Test equal block head/tail content",
			args: args{
				diffs: []DiffLine{
					{Type: Equal, Text: "1"},
					{Type: Equal, Text: "2"},
					{Type: Equal, Text: "3"},
					{Type: Delete, Text: "4"},
					{Type: Insert, Text: "5"},
					{Type: Insert, Text: "6"},
					{Type: Equal, Text: "7"},
					{Type: Equal, Text: "8"},
					{Type: Equal, Text: "9"},
					{Type: Equal, Text: "10"},
					{Type: Equal, Text: "11"},
					{Type: Delete, Text: "12"},
					{Type: Insert, Text: "13"},
					{Type: Equal, Text: "14"},
					{Type: Equal, Text: "15"},
					{Type: Equal, Text: "16"},
				},
				precontext:  2,
				postcontext: 2,
			},
			want: []Hunk{
				{
					Diffs: []DiffLine{
						{Type: Equal, Text: "2"},
						{Type: Equal, Text: "3"},
						{Type: Delete, Text: "4"},
						{Type: Insert, Text: "5"},
						{Type: Insert, Text: "6"},
						{Type: Equal, Text: "7"},
						{Type: Equal, Text: "8"},
					},
					SrcStart: 2,
					SrcLines: 5,
					DstStart: 2,
					DstLines: 6,
				},
				{
					Diffs: []DiffLine{
						{Type: Equal, Text: "10"},
						{Type: Equal, Text: "11"},
						{Type: Delete, Text: "12"},
						{Type: Insert, Text: "13"},
						{Type: Equal, Text: "14"},
						{Type: Equal, Text: "15"},
					},
					SrcStart: 8,
					SrcLines: 5,
					DstStart: 9,
					DstLines: 5,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeHunks(
				tt.args.diffs,
				tt.args.precontext,
				tt.args.postcontext,
			); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeHunks() = %v, want %v", got, tt.want)
			}
		})
	}
}
