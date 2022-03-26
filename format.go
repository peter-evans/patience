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

// UnifiedDiffOptions represents the options for UnifiedDiffTextWithOptions.
type UnifiedDiffOptions struct {
	// Precontext is the number of lines of context before each change in a hunk.
	Precontext int
	// Postcontext is the number of lines of context after each change in a hunk.
	Postcontext int
	// SrcHeader is the header for the source file.
	SrcHeader string
	// DstHeader is the header for the destination file.
	DstHeader string
}

// UnifiedDiffTextWithOptions returns the diff text in unidiff format.
func UnifiedDiffTextWithOptions(diffs []DiffLine, opts UnifiedDiffOptions) string {
	hunks := makeHunks(diffs, opts.Precontext, opts.Postcontext)
	s := []string{}
	if len(opts.SrcHeader) > 0 {
		s = append(s, fmt.Sprintf("--- %s", opts.SrcHeader))
	}
	if len(opts.DstHeader) > 0 {
		s = append(s, fmt.Sprintf("+++ %s", opts.DstHeader))
	}
	for _, h := range hunks {
		s = append(s, fmt.Sprintf("@@ -%d,%d +%d,%d @@", h.SrcStart, h.SrcLines, h.DstStart, h.DstLines))
		for _, l := range h.Diffs {
			if l.Type == Equal && len(l.Text) == 0 {
				s = append(s, "")
			} else {
				s = append(s, fmt.Sprintf("%s%s", typeSymbol(l.Type), l.Text))
			}
		}
	}
	return strings.Join(s, "\n")
}

// UnifiedDiffText returns the diff text in unidiff format with a context of 3 lines.
func UnifiedDiffText(diffs []DiffLine) string {
	return UnifiedDiffTextWithOptions(
		diffs,
		UnifiedDiffOptions{Precontext: 3, Postcontext: 3},
	)
}
