package state

// JSON Schema is generated from the Go structs above, never hand-written:
// two hand-maintained copies of one shape drift, and the drift is silent.
// The generated files live in shared/schemas/pipeline/ because agents (and
// humans) read them from the installed framework content.

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
)

// SchemaDir is where generated pipeline schemas are committed, relative to
// the repo root.
const SchemaDir = "shared/schemas/pipeline"

// StageSchema names one stage's state document and the type behind it.
type StageSchema struct {
	// Stage is the stage ID this schema belongs to, matching the plan's
	// stable stage IDs (never an ordinal).
	Stage string
	// FileName is the committed schema file, e.g. analysis.schema.json.
	FileName string
	subject  interface{}
}

// StageSchemas returns every stage whose state is typed today. L2.9's first
// cut types one hop; the rest of the pipeline still exchanges markdown, and
// this list is where later epics grow.
func StageSchemas() []StageSchema {
	return []StageSchema{
		{Stage: "analyst", FileName: "analysis.schema.json", subject: &AnalysisState{}},
		{Stage: "architect", FileName: "architecture.schema.json", subject: &ArchitectureState{}},
	}
}

// SchemaForStage returns the generated JSON Schema for a stage, or false
// when that stage is not typed yet.
func SchemaForStage(stageID string) ([]byte, bool) {
	for _, schema := range StageSchemas() {
		if schema.Stage != stageID {
			continue
		}
		raw, err := schema.Generate()
		if err != nil {
			return nil, false
		}
		return raw, true
	}
	return nil, false
}

// Generate renders the JSON Schema for one stage. Output is deterministic:
// the same structs always produce byte-identical schemas, which is what
// makes the committed copies checkable.
func (s StageSchema) Generate() ([]byte, error) {
	reflector := &jsonschema.Reflector{ExpandedStruct: true, DoNotReference: true}
	schema := reflector.Reflect(s.subject)
	var rendered bytes.Buffer
	encoder := json.NewEncoder(&rendered)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(schema); err != nil {
		return nil, fmt.Errorf("render schema for stage %q: %w", s.Stage, err)
	}
	return rendered.Bytes(), nil
}
