// Package patience implements the Patience Diff algorithm.
package patience

import (
	"fmt"
	"strings"
)

// typeSymbol returns the associated symbol of a DiffType.
func typeSymbol(t DiffType) string {
	switch t {
	case Equal:
		return " "
	case Insert:
		return "+"
	case Delete:
		return "-"
	default:
		panic("unknown DiffType")
	}
}

// DiffText returns the source and destination texts (all equalities, insertions and deletions).
func DiffText(diffs []DiffLine) string {
	s := make([]string, len(diffs))
	for i, l := range diffs {
		if len(l.Text) == 0 && l.Type == Equal {
			continue
		}
		s[i] = fmt.Sprintf("%s%s", typeSymbol(l.Type), l.Text)
	}
	return strings.Join(s, "\n")
}

// DiffTextA returns the source text (all equalities and deletions).
func DiffTextA(diffs []DiffLine) string {
	s := []string{}
	for _, l := range diffs {
		if l.Type == Insert {
			continue
		}
		if l.Type == Equal && len(l.Text) == 0 {
			s = append(s, "")
		} else {
			s = append(s, fmt.Sprintf("%s%s", typeSymbol(l.Type), l.Text))
		}
	}
	return strings.Join(s, "\n")
}

// DiffTextB returns the destination text (all equalities and insertions).
func DiffTextB(diffs []DiffLine) string {
	s := []string{}
	for _, l := range diffs {
		if l.Type == Delete {
			continue
		}
		if l.Type == Equal && len(l.Text) == 0 {
			s = append(s, "")
		} else {
			s = append(s, fmt.Sprintf("%s%s", typeSymbol(l.Type), l.Text))
		}
	}
	return strings.Join(s, "\n")
}
