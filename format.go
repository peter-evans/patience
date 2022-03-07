// Package patience implements the Patience Diff algorithm.
package patience

import (
	"fmt"
	"strings"
)

// DiffText returns the source and destination texts (all equalities, insertions and deletions).
func DiffText(diffs []DiffLine) string {
	s := make([]string, len(diffs))
	for i, l := range diffs {
		prefix := " "
		switch l.Type {
		case Insert:
			prefix = "+"
		case Delete:
			prefix = "-"
		}
		s[i] = fmt.Sprintf("%s%s", prefix, l.Text)
	}
	return strings.Join(s, "\n")
}
