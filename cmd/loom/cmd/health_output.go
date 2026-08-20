package cmd

import (
	"fmt"
	"io"
)

type healthOutput struct{ writer io.Writer }

func (output healthOutput) heading(target string) {
	fmt.Fprintf(output.writer, "loom health: checking %s\n\n", target)
}

func (output healthOutput) success(message string) {
	fmt.Fprintln(output.writer, "  ✓ "+message)
}

func (output healthOutput) warning(message string) {
	fmt.Fprintln(output.writer, "  ⚠ "+message)
}

func (output healthOutput) failure(message string) {
	fmt.Fprintln(output.writer, "  ✗ "+message)
}

func (output healthOutput) note() {
	fmt.Fprintln(output.writer, "\nnote: for full framework health checks, run scripts/health-check.sh in the framework repo")
}

func (output healthOutput) summary(failures, warnings int) {
	if failures == 0 && warnings == 0 {
		fmt.Fprintln(output.writer, "\nAll checks passed.")
		return
	}
	result := countLabel(failures, "failure")
	if failures == 0 {
		result = countLabel(warnings, "warning")
	} else if warnings > 0 {
		result += ", " + countLabel(warnings, "warning")
	}
	fmt.Fprintf(output.writer, "\n%s. Run `loom install` to repair.\n", result)
}

func countLabel(count int, singular string) string {
	label := singular
	if count != 1 {
		label += "s"
	}
	return fmt.Sprintf("%d %s", count, label)
}
