// Package patience implements the Patience Diff algorithm.
package patience

// LCS computes the longest common subsequence of two string
// slices and returns the index pairs of the LCS.
func LCS(a, b []string) [][2]int {
	// Initialize the LCS table.
	lcs := make([][]int, len(a)+1)
	for i := 0; i <= len(a); i++ {
		lcs[i] = make([]int, len(b)+1)
	}

	// Populate the LCS table.
	for i := 1; i < len(lcs); i++ {
		for j := 1; j < len(lcs[i]); j++ {
			if a[i-1] == b[j-1] {
				lcs[i][j] = lcs[i-1][j-1] + 1
			} else {
				lcs[i][j] = max(lcs[i-1][j], lcs[i][j-1])
			}
		}
	}

	// Backtrack to find the LCS.
	i, j := len(a), len(b)
	s := make([][2]int, 0, lcs[i][j])
	for i > 0 && j > 0 {
		switch {
		case a[i-1] == b[j-1]:
			s = append(s, [2]int{i - 1, j - 1})
			i--
			j--
		case lcs[i-1][j] > lcs[i][j-1]:
			i--
		default:
			j--
		}
	}

	// Reverse the backtracked LCS.
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

// max returns the maximum of two integers.
// nolint:predeclared
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
