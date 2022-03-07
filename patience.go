// Package patience implements the Patience Diff algorithm.
package patience

// DiffType defines the type of a diff element.
type DiffType int8

const (
	// Delete represents a diff delete operation.
	Delete DiffType = -1
	// Insert represents a diff insert operation.
	Insert DiffType = 1
	// Equal represents no diff.
	Equal DiffType = 0
)

// DiffLine represents a single line and its diff type.
type DiffLine struct {
	Text string
	Type DiffType
}

// toDiffLines is a convenience function to convert a slice of strings
// to a slice of DiffLines with the specified diff type.
func toDiffLines(a []string, t DiffType) []DiffLine {
	diffs := make([]DiffLine, len(a))
	for i, l := range a {
		diffs[i] = DiffLine{l, t}
	}
	return diffs
}

// uniqueElements returns a slice of unique elements from a slice of
// strings, and a slice of the original indices of each element.
func uniqueElements(a []string) ([]string, []int) {
	m := make(map[string]int)
	for _, e := range a {
		m[e]++
	}
	elements := []string{}
	indices := []int{}
	for i, e := range a {
		if m[e] == 1 {
			elements = append(elements, e)
			indices = append(indices, i)
		}
	}
	return elements, indices
}

// Diff returns the patience diff of two slices of strings.
func Diff(a, b []string) []DiffLine {
	switch {
	case len(a) == 0 && len(b) == 0:
		return nil
	case len(a) == 0:
		return toDiffLines(b, Insert)
	case len(b) == 0:
		return toDiffLines(a, Delete)
	}

	// Find equal elements at the head of slices a and b.
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	if i > 0 {
		return append(
			toDiffLines(a[:i], Equal),
			Diff(a[i:], b[i:])...,
		)
	}

	// Find equal elements at the tail of slices a and b.
	j := 0
	for j < len(a) && j < len(b) && a[len(a)-1-j] == b[len(b)-1-j] {
		j++
	}
	if j > 0 {
		return append(
			Diff(a[:len(a)-j], b[:len(b)-j]),
			toDiffLines(a[len(a)-j:], Equal)...,
		)
	}

	// Find the longest common subsequence of unique elements in a and b.
	ua, idxa := uniqueElements(a)
	ub, idxb := uniqueElements(b)
	lcs := LCS(ua, ub)

	// If the LCS is empty, the diff is all deletions and insertions.
	if len(lcs) == 0 {
		return append(toDiffLines(a, Delete), toDiffLines(b, Insert)...)
	}

	// Lookup the original indices of slices a and b.
	for i, x := range lcs {
		lcs[i][0] = idxa[x[0]]
		lcs[i][1] = idxb[x[1]]
	}

	diffs := []DiffLine{}
	ga, gb := 0, 0
	for _, ip := range lcs {
		// Diff the gaps between the lcs elements.
		diffs = append(diffs, Diff(a[ga:ip[0]], b[gb:ip[1]])...)
		// Append the LCS elements to the diff.
		diffs = append(diffs, DiffLine{Type: Equal, Text: a[ip[0]]})
		ga = ip[0] + 1
		gb = ip[1] + 1
	}
	// Diff the remaining elements of a and b after the final LCS element.
	diffs = append(diffs, Diff(a[ga:], b[gb:])...)

	return diffs
}
