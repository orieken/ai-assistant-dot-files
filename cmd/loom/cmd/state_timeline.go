package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/spf13/cobra"
)

func runStateTimeline(cmd *cobra.Command, _ []string) error {
	timeline, err := openTimeline(stateArgs.spec)
	if err != nil {
		return err
	}
	events, err := timeline.Read()
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return fmt.Errorf("no events recorded for %s yet", stateArgs.spec)
	}
	if stateArgs.asJSON {
		return printEventsJSON(cmd, events)
	}
	printEventReport(cmd, events)
	return nil
}

func printEventsJSON(cmd *cobra.Command, events []orchestrator.Event) error {
	raw, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("encode events: %w", err)
	}
	cmd.Println(string(raw))
	return nil
}

// printEventReport renders the timeline with elapsed time from the first
// event — measured by subtracting timestamps this binary wrote, not
// estimated by anything.
func printEventReport(cmd *cobra.Command, events []orchestrator.Event) {
	start := events[0].At
	for _, event := range events {
		cmd.Printf("%8s  %-20s %s\n",
			elapsed(start, event.At), event.Kind, strings.TrimSpace(subjectOf(event)))
	}
}

func elapsed(start, at time.Time) string {
	return at.Sub(start).Round(time.Second).String()
}

func subjectOf(event orchestrator.Event) string {
	if event.Kind == orchestrator.EventArtifactCorrected {
		return correctionSubject(event)
	}
	parts := []string{event.Stage, event.Gate, string(event.StaleReason), string(event.ApprovalMethod), event.Error}
	kept := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			kept = append(kept, part)
		}
	}
	return strings.Join(kept, " ")
}

// correctionSubject reads as "who was corrected, by how much" rather than
// dumping an absolute diff path onto the line. The path is in --json for
// anything that wants to open the diff.
func correctionSubject(event orchestrator.Event) string {
	return strings.TrimSpace(fmt.Sprintf("%s (%s) %s", event.Stage, event.Agent, event.Correction))
}
