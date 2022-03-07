// Package patience implements the Patience Diff algorithm.
package patience

import (
	"reflect"
	"strings"
	"testing"
)

func Test_uniqueElements(t *testing.T) {
	type args struct {
		a []string
	}
	tests := []struct {
		name        string
		args        args
		wantOut     []string
		wantIndices []int
	}{
		{
			name: "Test every element is unique",
			args: args{
				a: []string{"a", "b", "c"},
			},
			wantOut:     []string{"a", "b", "c"},
			wantIndices: []int{0, 1, 2},
		},
		{
			name: "Test duplicate elements",
			args: args{
				a: []string{"a", "b", "a", "c"},
			},
			wantOut:     []string{"b", "c"},
			wantIndices: []int{1, 3},
		},
		{
			name: "Test no unique elements",
			args: args{
				a: []string{"a", "b", "a", "c", "c", "b"},
			},
			wantOut:     []string{},
			wantIndices: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, gotIndices := uniqueElements(tt.args.a)
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("uniqueElements() = %v, want %v", gotOut, tt.wantOut)
			}
			if !reflect.DeepEqual(gotIndices, tt.wantIndices) {
				t.Errorf("uniqueElements() = %v, want %v", gotIndices, tt.wantIndices)
			}
		})
	}
}

func TestDiff(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []DiffLine
	}{
		{
			name: "Test empty slices",
			args: args{
				a: []string{},
				b: []string{},
			},
			want: nil,
		},
		{
			name: "Test empty slice a",
			args: args{
				a: []string{},
				b: []string{"a"},
			},
			want: []DiffLine{
				{Text: "a", Type: Insert},
			},
		},
		{
			name: "Test empty slice b",
			args: args{
				a: []string{"a"},
				b: []string{},
			},
			want: []DiffLine{
				{Text: "a", Type: Delete},
			},
		},
		{
			name: "Test no diff",
			args: args{
				a: []string{"a"},
				b: []string{"a"},
			},
			want: []DiffLine{
				{Text: "a", Type: Equal},
			},
		},
		{
			name: "Test equal elements at the head of slices a and b",
			args: args{
				a: []string{"a", "b"},
				b: []string{"a", "c"},
			},
			want: []DiffLine{
				{Text: "a", Type: Equal},
				{Text: "b", Type: Delete},
				{Text: "c", Type: Insert},
			},
		},
		{
			name: "Test equal elements at the tail of slices a and b",
			args: args{
				a: []string{"a", "c"},
				b: []string{"b", "c"},
			},
			want: []DiffLine{
				{Text: "a", Type: Delete},
				{Text: "b", Type: Insert},
				{Text: "c", Type: Equal},
			},
		},
		{
			name: "Test equal elements at the head and tail of slices a and b",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"a", "d", "c"},
			},
			want: []DiffLine{
				{Text: "a", Type: Equal},
				{Text: "b", Type: Delete},
				{Text: "d", Type: Insert},
				{Text: "c", Type: Equal},
			},
		},
		{
			name: "Test diffing the gaps between the LCS elements of slices a and b",
			args: args{
				a: []string{"a", "w", "b", "x", "c"},
				b: []string{"a", "y", "b", "z", "c"},
			},
			want: []DiffLine{
				{Text: "a", Type: Equal},
				{Text: "w", Type: Delete},
				{Text: "y", Type: Insert},
				{Text: "b", Type: Equal},
				{Text: "x", Type: Delete},
				{Text: "z", Type: Insert},
				{Text: "c", Type: Equal},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Diff(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestPatienceCanonical tests the "canonical" patience diff example.
// https://alfedenzo.livejournal.com/170301.html
func TestPatienceCanonical(t *testing.T) {
	a := strings.Split(`#include <stdio.h>

// Frobs foo heartily
int frobnitz(int foo)
{
    int i;
    for(i = 0; i < 10; i++)
    {
        printf("Your answer is: ");
        printf("%d\n", foo);
    }
}

int fact(int n)
{
    if(n > 1)
    {
        return fact(n-1) * n;
    }
    return 1;
}

int main(int argc, char **argv)
{
    frobnitz(fact(10));
}`, "\n")

	b := strings.Split(`#include <stdio.h>

int fib(int n)
{
    if(n > 2)
    {
        return fib(n-1) + fib(n-2);
    }
    return 1;
}

// Frobs foo heartily
int frobnitz(int foo)
{
    int i;
    for(i = 0; i < 10; i++)
    {
        printf("%d\n", foo);
    }
}

int main(int argc, char **argv)
{
    frobnitz(fib(10));
}`, "\n")

	want := ` #include <stdio.h>

+int fib(int n)
+{
+    if(n > 2)
+    {
+        return fib(n-1) + fib(n-2);
+    }
+    return 1;
+}
+
 // Frobs foo heartily
 int frobnitz(int foo)
 {
     int i;
     for(i = 0; i < 10; i++)
     {
-        printf("Your answer is: ");
         printf("%d\n", foo);
     }
 }

-int fact(int n)
-{
-    if(n > 1)
-    {
-        return fact(n-1) * n;
-    }
-    return 1;
-}
-
 int main(int argc, char **argv)
 {
-    frobnitz(fact(10));
+    frobnitz(fib(10));
 }`

	if got := DiffText(Diff(a, b)); got != want {
		t.Errorf("TestPatienceCanonical = %v, want %v", got, want)
	}
}
