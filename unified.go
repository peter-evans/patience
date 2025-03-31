package patience

// Hunk represents a subsection of a diff.
type Hunk struct {
	Diffs    []DiffLine
	SrcStart int
	SrcLines int
	DstStart int
	DstLines int
}

// makeHunks returns the hunks of a diff.
func makeHunks(diffs []DiffLine, precontext, postcontext int) []Hunk {
	if len(diffs) == 0 {
		return nil
	}

	hunks := []Hunk{}

	// Update hunks with a diff block.
	updateHunks := func(block Hunk, lastBlock bool) {
		curHunk := len(hunks) - 1
		if block.Diffs[0].Type == Equal {
			// Unmodified block.
			if len(hunks) == 0 {
				// Start a new hunk with the tail of the block.
				ctxLen := min(precontext, len(block.Diffs))
				hunks = append(
					hunks,
					Hunk{
						Diffs:    block.Diffs[len(block.Diffs)-ctxLen:],
						SrcStart: len(block.Diffs) - ctxLen + block.SrcStart,
						SrcLines: ctxLen,
						DstStart: len(block.Diffs) - ctxLen + block.DstStart,
						DstLines: ctxLen,
					},
				)
			} else {
				// Update the current hunk.
				maxNonContext := precontext + postcontext
				if lastBlock {
					maxNonContext = postcontext
				}
				if len(block.Diffs) <= maxNonContext {
					// Block is small enough to be appended to the current hunk.
					hunks[curHunk].Diffs = append(hunks[curHunk].Diffs, block.Diffs...)
					hunks[curHunk].SrcLines += len(block.Diffs)
					hunks[curHunk].DstLines += len(block.Diffs)
				} else {
					// Append the head of the block to the current hunk.
					hunks[curHunk].Diffs = append(hunks[curHunk].Diffs, block.Diffs[:postcontext]...)
					hunks[curHunk].SrcLines += postcontext
					hunks[curHunk].DstLines += postcontext
					if !lastBlock {
						// Start a new hunk with the tail of the block.
						hunks = append(
							hunks,
							Hunk{
								Diffs:    block.Diffs[len(block.Diffs)-precontext:],
								SrcStart: len(block.Diffs) - precontext + block.SrcStart,
								SrcLines: precontext,
								DstStart: len(block.Diffs) - precontext + block.DstStart,
								DstLines: precontext,
							},
						)
					}
				}
				// Update starting line numbers if the current hunk had no source or destination diff.
				if hunks[curHunk].SrcStart == 0 {
					hunks[curHunk].SrcStart = block.SrcStart
				}
				if hunks[curHunk].DstStart == 0 {
					hunks[curHunk].DstStart = block.DstStart
				}
			}
		} else {
			// Modified block.
			if len(hunks) > 0 {
				hunks[curHunk].Diffs = append(hunks[curHunk].Diffs, block.Diffs...)
				hunks[curHunk].SrcLines += block.SrcLines
				hunks[curHunk].DstLines += block.DstLines
			} else {
				hunks = append(
					hunks,
					Hunk{
						Diffs:    block.Diffs,
						SrcStart: block.SrcStart,
						SrcLines: block.SrcLines,
						DstStart: block.DstStart,
						DstLines: block.DstLines,
					},
				)
			}
		}
	}

	// Aggregate blocks of modified and unmodified diff lines, creating
	// or updating hunks after each block.
	var block Hunk
	modifiedLines := 0
	srcLineNum, dstLineNum := 0, 0
	for _, l := range diffs {
		if len(block.Diffs) == 0 ||
			block.Diffs[0].Type == l.Type ||
			(block.Diffs[0].Type != l.Type && block.Diffs[0].Type != Equal && l.Type != Equal) {
			block.Diffs = append(block.Diffs, l)
		} else {
			updateHunks(block, false)
			block = Hunk{Diffs: []DiffLine{l}}
		}

		switch l.Type {
		case Delete:
			srcLineNum++
			block.SrcLines++
			modifiedLines++
		case Insert:
			dstLineNum++
			block.DstLines++
			modifiedLines++
		case Equal:
			srcLineNum++
			dstLineNum++
			block.SrcLines++
			block.DstLines++
		}

		if block.SrcStart == 0 && (l.Type == Equal || l.Type == Delete) {
			block.SrcStart = srcLineNum
		}
		if block.DstStart == 0 && (l.Type == Equal || l.Type == Insert) {
			block.DstStart = dstLineNum
		}
	}
	updateHunks(block, true)

	// Return no hunks if the diffs contain only equal lines.
	if modifiedLines == 0 {
		return nil
	}

	return hunks
}

// min returns the minimum of two integers.
// nolint:predeclared
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
