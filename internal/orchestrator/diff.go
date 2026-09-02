package orchestrator

// A minimal unified diff, used to describe what a human changed in an
// artifact (roadmap L4.5). The standard library has none, and pulling a
// dependency in for one call site is not worth the supply-chain surface
// when the algorithm is textbook.
//
// Line-level LCS, which is what `diff` itself uses. Artifacts are hundreds
// of lines, so the O(n*m) table is a few hundred kilobytes at worst — not
// worth the complexity of Myers for this input size.

import (
	"fmt"
	"strings"
)

// DiffStat counts what changed, for a timeline event that must stay one
// greppable line while the detail lives on disk.
type DiffStat struct {
	Added   int `json:"added"`
	Removed int `json:"removed"`
}

// Changed reports whether anything differs at all.
func (d DiffStat) Changed() bool { return d.Added > 0 || d.Removed > 0 }

func (d DiffStat) String() string { return fmt.Sprintf("+%d/-%d", d.Added, d.Removed) }

// UnifiedDiff renders before → after in unified format and counts the
// changed lines. An unchanged pair yields an empty diff and a zero stat, so
// callers test the stat rather than the string.
func UnifiedDiff(before, after string) (string, DiffStat) {
	beforeLines, afterLines := splitLines(before), splitLines(after)
	edits := lineEdits(beforeLines, afterLines)
	return renderEdits(edits), statOf(edits)
}

// editKind is one line's fate in the diff.
type editKind byte

const (
	editKeep editKind = ' '
	editAdd  editKind = '+'
	editDrop editKind = '-'
)

type lineEdit struct {
	kind editKind
	text string
}

func splitLines(text string) []string {
	if text == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(text, "\n"), "\n")
}

// lineEdits walks the LCS table backwards to produce the edit script.
func lineEdits(before, after []string) []lineEdit {
	table := lcsTable(before, after)
	edits := make([]lineEdit, 0, len(before)+len(after))
	row, column := 0, 0
	for row < len(before) && column < len(after) {
		edits, row, column = stepThrough(edits, before, after, table, row, column)
	}
	for ; row < len(before); row++ {
		edits = append(edits, lineEdit{editDrop, before[row]})
	}
	for ; column < len(after); column++ {
		edits = append(edits, lineEdit{editAdd, after[column]})
	}
	return edits
}

// stepThrough advances one cell of the LCS table, appending the edit that
// cell implies.
func stepThrough(edits []lineEdit, before, after []string, table [][]int, row, column int) ([]lineEdit, int, int) {
	switch {
	case before[row] == after[column]:
		return append(edits, lineEdit{editKeep, before[row]}), row + 1, column + 1
	case table[row+1][column] >= table[row][column+1]:
		return append(edits, lineEdit{editDrop, before[row]}), row + 1, column
	default:
		return append(edits, lineEdit{editAdd, after[column]}), row, column + 1
	}
}

// lcsTable[i][j] is the longest common subsequence length of before[i:] and
// after[j:].
func lcsTable(before, after []string) [][]int {
	table := make([][]int, len(before)+1)
	for row := range table {
		table[row] = make([]int, len(after)+1)
	}
	for row := len(before) - 1; row >= 0; row-- {
		for column := len(after) - 1; column >= 0; column-- {
			table[row][column] = cellValue(before, after, table, row, column)
		}
	}
	return table
}

func cellValue(before, after []string, table [][]int, row, column int) int {
	if before[row] == after[column] {
		return table[row+1][column+1] + 1
	}
	return max(table[row+1][column], table[row][column+1])
}

func statOf(edits []lineEdit) DiffStat {
	var stat DiffStat
	for _, edit := range edits {
		switch edit.kind {
		case editAdd:
			stat.Added++
		case editDrop:
			stat.Removed++
		case editKeep:
		}
	}
	return stat
}

// renderEdits emits the whole file rather than hunks with context windows.
// Artifacts are small, a whole-file diff is unambiguous, and hunk headers
// are the part of unified format that is easy to get subtly wrong.
func renderEdits(edits []lineEdit) string {
	if !statOf(edits).Changed() {
		return ""
	}
	var builder strings.Builder
	for _, edit := range edits {
		builder.WriteByte(byte(edit.kind))
		builder.WriteString(edit.text)
		builder.WriteByte('\n')
	}
	return builder.String()
}
