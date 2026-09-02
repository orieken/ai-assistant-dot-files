package orchestrator_test

import (
	"strings"
	"testing"

	"github.com/orieken/loom/internal/orchestrator"
)

func TestUnifiedDiff(t *testing.T) {
	cases := []struct {
		name          string
		before, after string
		added         int
		removed       int
		wantLines     []string
	}{
		{
			name:   "identical text is not a change",
			before: "one\ntwo\n", after: "one\ntwo\n",
		},
		{
			name:   "an appended line",
			before: "one\n", after: "one\ntwo\n",
			added: 1, wantLines: []string{" one", "+two"},
		},
		{
			name:   "a removed line",
			before: "one\ntwo\n", after: "one\n",
			removed: 1, wantLines: []string{" one", "-two"},
		},
		{
			name:   "a replaced line counts as both",
			before: "one\ntwo\n", after: "one\nTWO\n",
			added: 1, removed: 1, wantLines: []string{" one", "-two", "+TWO"},
		},
		{
			name:   "an insertion in the middle keeps surrounding lines",
			before: "a\nc\n", after: "a\nb\nc\n",
			added: 1, wantLines: []string{" a", "+b", " c"},
		},
		{
			name:   "empty to content",
			before: "", after: "one\n",
			added: 1, wantLines: []string{"+one"},
		},
		{
			name:   "content to empty",
			before: "one\n", after: "",
			removed: 1, wantLines: []string{"-one"},
		},
		{
			// A trailing newline is not a change: artifacts get rewritten by
			// editors that add or drop one, and reporting that as a human
			// correction would be noise indistinguishable from signal.
			name:   "a trailing newline alone is not a change",
			before: "one\ntwo\n", after: "one\ntwo",
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			diff, stat := orchestrator.UnifiedDiff(testCase.before, testCase.after)
			if stat.Added != testCase.added || stat.Removed != testCase.removed {
				t.Errorf("stat = %v, want +%d/-%d", stat, testCase.added, testCase.removed)
			}
			if stat.Changed() != (testCase.added+testCase.removed > 0) {
				t.Errorf("Changed() = %v for stat %v", stat.Changed(), stat)
			}
			assertDiffLines(t, diff, testCase.wantLines)
		})
	}
}

func assertDiffLines(t *testing.T, diff string, want []string) {
	t.Helper()
	if len(want) == 0 {
		if diff != "" {
			t.Errorf("diff = %q, want empty when nothing changed", diff)
		}
		return
	}
	got := strings.Split(strings.TrimSuffix(diff, "\n"), "\n")
	if len(got) != len(want) {
		t.Fatalf("diff =\n%s\nwant lines %v", diff, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("diff line %d = %q, want %q\nfull diff:\n%s", i, got[i], want[i], diff)
		}
	}
}

func TestDiffStatStringIsCompactForATimelineLine(t *testing.T) {
	_, stat := orchestrator.UnifiedDiff("a\nb\n", "a\nB\nc\n")
	if stat.String() != "+2/-1" {
		t.Errorf("String() = %q, want %q", stat.String(), "+2/-1")
	}
}
