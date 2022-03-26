// Package patience implements the Patience Diff algorithm.
package patience

import (
	"strings"
	"testing"
)

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

// TestPatienceUnified tests the "canonical" patience diff example formatted as unidiff.
func TestPatienceUnified(t *testing.T) {
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

	want := `@@ -1,26 +1,25 @@
 #include <stdio.h>

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

	if got := UnifiedDiffText(Diff(a, b)); got != want {
		t.Errorf("TestPatienceUnified = %v, want %v", got, want)
	}
}

func TestPlainDiffReadmeExample(t *testing.T) {
	a := strings.Split(`the
quick
brown
chicken
jumps
over
the
dog`, "\n")

	b := strings.Split(`the
quick
brown
fox
jumps
over
the
lazy
dog`, "\n")

	want := ` the
 quick
 brown
-chicken
+fox
 jumps
 over
 the
+lazy
 dog`

	if got := DiffText(Diff(a, b)); got != want {
		t.Errorf("TestPlainDiffReadmeExample = %v, want %v", got, want)
	}
}

func TestUnifiedDiffReadmeExample(t *testing.T) {
	a := strings.Split(`the
quick
brown
chicken
jumps
over
the
dog`, "\n")

	b := strings.Split(`the
quick
brown
fox
jumps
over
the
lazy
dog`, "\n")

	want := `--- a.txt
+++ b.txt
@@ -3,3 +3,3 @@
 brown
-chicken
+fox
 jumps
@@ -7,2 +7,3 @@
 the
+lazy
 dog`

	if got := UnifiedDiffTextWithOptions(
		Diff(a, b),
		UnifiedDiffOptions{Precontext: 1, Postcontext: 1, SrcHeader: "a.txt", DstHeader: "b.txt"},
	); got != want {
		t.Errorf("TestUnifiedDiffReadmeExample = %v, want %v", got, want)
	}
}
