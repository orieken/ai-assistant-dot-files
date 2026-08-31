package claude

// A typed stage (roadmap L2.9) returns a state document instead of
// markdown. The agent definition is not edited to say so — those files are
// shared with the markdown pipeline, which must keep working — so the
// schema and the output instruction are appended here, at invocation time.

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/orieken/loom/internal/orchestrator"
	"github.com/orieken/loom/internal/state"
)

// typedInstruction appends the schema and the output contract for a typed
// stage. The schema is inlined rather than referenced by path so the run
// does not depend on the framework being installed in the target project.
func typedInstruction(stage orchestrator.Stage, input orchestrator.StageInput) (string, error) {
	schema, ok := state.SchemaForKind(state.Kind(stage.StateKind))
	if !ok {
		return "", fmt.Errorf("stage %q declares state kind %q, which has no schema", stage.ID, stage.StateKind)
	}
	var instruction strings.Builder
	instruction.WriteString("\n---\n\nOUTPUT CONTRACT (this overrides any output-format instruction above).\n")
	instruction.WriteString("Return a single JSON object conforming to this schema, and nothing else.\n")
	instruction.WriteString("Do not write files. Do not add commentary before or after the JSON.\n\n")
	instruction.Write(schema)
	instruction.WriteString(upstreamSection(input))
	return instruction.String(), nil
}

// upstreamSection hands the agent the projected fields of the previous
// stage — the data it is allowed to read, rather than a document to reparse.
func upstreamSection(input orchestrator.StageInput) string {
	if len(input.UpstreamState) == 0 {
		return ""
	}
	return fmt.Sprintf("\n\nInput from the %s stage (these fields are your source of truth; do not go looking for its markdown):\n\n%s\n",
		input.UpstreamStage, input.UpstreamState)
}

// extractJSON pulls the state document out of an agent's response. Raw JSON
// is accepted, and so is exactly one fenced block with only whitespace
// around it — models fence their output as a formatting habit, and failing
// a run over that would report a reflex as a modelling error. Anything
// else fails: scanning for the first balanced object anywhere would happily
// accept a schema example the agent quoted back.
func extractJSON(response []byte) ([]byte, error) {
	trimmed := bytes.TrimSpace(response)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("agent returned nothing")
	}
	if trimmed[0] == '{' {
		return trimmed, nil
	}
	return unfence(trimmed)
}

func unfence(trimmed []byte) ([]byte, error) {
	text := string(trimmed)
	if !strings.HasPrefix(text, "```") || !strings.HasSuffix(text, "```") {
		return nil, unexpectedResponse(text)
	}
	body := strings.TrimSuffix(text, "```")
	body = body[strings.Index(body, "\n")+1:]
	body = strings.TrimSpace(body)
	if !strings.HasPrefix(body, "{") || strings.Contains(body, "```") {
		return nil, unexpectedResponse(text)
	}
	return []byte(body), nil
}

func unexpectedResponse(text string) error {
	return fmt.Errorf("agent did not return a JSON state document — got: %s", truncate(text, 800))
}
