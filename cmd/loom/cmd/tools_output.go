package cmd

import (
	"fmt"
	"io"
)

type toolsOutput struct{ writer io.Writer }

func (o toolsOutput) heading(title string) {
	fmt.Fprintln(o.writer, title)
	fmt.Fprintln(o.writer)
}

func (o toolsOutput) sectionHeading(tier toolTier) {
	switch tier {
	case tierHigh:
		fmt.Fprintln(o.writer, "  tier: high  (direct pipeline integration)")
	case tierMedium:
		fmt.Fprintln(o.writer, "  tier: medium  (opt-in, adopt when the pain point arises)")
	}
}

func (o toolsOutput) toolRow(name, status, description string) {
	fmt.Fprintf(o.writer, "  %s %-16s %s\n", status, name, description)
}

func (o toolsOutput) installNote(note string) {
	fmt.Fprintf(o.writer, "              install: %s\n", note)
}

func (o toolsOutput) postNote(note string) {
	fmt.Fprintf(o.writer, "              note:    %s\n", note)
}

func (o toolsOutput) separator() {
	fmt.Fprintln(o.writer)
}

func (o toolsOutput) statusSummary(missing int) {
	if missing == 0 {
		fmt.Fprintln(o.writer, "All context tools are installed.")
		return
	}
	noun := "tool"
	if missing != 1 {
		noun = "tools"
	}
	fmt.Fprintf(o.writer, "%d %s not installed. Run `loom tools install --tier high` to install the high-priority ones.\n", missing, noun)
}

func (o toolsOutput) installHeading(count int) {
	noun := "tool"
	if count != 1 {
		noun = "tools"
	}
	fmt.Fprintf(o.writer, "installing %d %s\n\n", count, noun)
}

func (o toolsOutput) installSkip(name string) {
	fmt.Fprintf(o.writer, "  ✓ %-16s already installed\n", name)
}

func (o toolsOutput) installBuiltIn(name string) {
	fmt.Fprintf(o.writer, "  ✓ %-16s built-in (no install needed)\n", name)
}

func (o toolsOutput) installManual(name, cmd string) {
	fmt.Fprintf(o.writer, "  ! %-16s manual install required:\n", name)
	fmt.Fprintf(o.writer, "              %s\n", cmd)
}

func (o toolsOutput) installSuccess(name string) {
	fmt.Fprintf(o.writer, "  ✓ %-16s installed\n", name)
}

func (o toolsOutput) installFailed(name string, err error) {
	fmt.Fprintf(o.writer, "  ✗ %-16s install failed: %v\n", name, err)
}

func (o toolsOutput) installPost(note string) {
	fmt.Fprintf(o.writer, "              next:    %s\n", note)
}

func (o toolsOutput) installSummary(installed, skipped, failed, manual int) {
	fmt.Fprintln(o.writer)
	if installed > 0 {
		fmt.Fprintf(o.writer, "installed: %d\n", installed)
	}
	if skipped > 0 {
		fmt.Fprintf(o.writer, "already installed: %d\n", skipped)
	}
	if manual > 0 {
		fmt.Fprintf(o.writer, "manual install required: %d (commands printed above)\n", manual)
	}
	if failed > 0 {
		fmt.Fprintf(o.writer, "failed: %d\n", failed)
	}
}
