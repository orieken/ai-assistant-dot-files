package orchestrator_test

// The fitness function L3.9 exists to create: adding an event type without
// documenting it must fail the build.
//
// The check reads the constants out of the source with go/ast rather than
// trusting a hand-kept list, because a hand-kept list is exactly what the
// item is replacing. A registry that enumerated itself would pass happily
// while a new constant sat undocumented beside it.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/orieken/loom/internal/orchestrator"
)

const vocabularySourceFile = "timeline.go"

// declaredEventKinds parses the source for every constant declared with type
// EventKind and returns its string value.
func declaredEventKinds(t *testing.T) []string {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), vocabularySourceFile, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", vocabularySourceFile, err)
	}
	kinds := make([]string, 0, 16)
	ast.Inspect(file, func(node ast.Node) bool {
		spec, ok := node.(*ast.ValueSpec)
		if !ok || !isEventKindSpec(spec) {
			return true
		}
		kinds = append(kinds, constantValues(t, spec)...)
		return true
	})
	return kinds
}

func isEventKindSpec(spec *ast.ValueSpec) bool {
	identifier, ok := spec.Type.(*ast.Ident)
	return ok && identifier.Name == "EventKind"
}

func constantValues(t *testing.T, spec *ast.ValueSpec) []string {
	t.Helper()
	values := make([]string, 0, len(spec.Values))
	for _, value := range spec.Values {
		literal, ok := value.(*ast.BasicLit)
		if !ok || literal.Kind != token.STRING {
			continue
		}
		unquoted, err := strconv.Unquote(literal.Value)
		if err != nil {
			t.Fatalf("unquote %s: %v", literal.Value, err)
		}
		values = append(values, unquoted)
	}
	return values
}

func TestEveryEventKindIsDocumented(t *testing.T) {
	documented := make(map[string]bool)
	for _, doc := range orchestrator.EventVocabulary() {
		documented[string(doc.Kind)] = true
	}
	for _, kind := range declaredEventKinds(t) {
		if !documented[kind] {
			t.Errorf("event kind %q is declared but not in EventVocabulary — add an EventDoc for it in "+
				"vocabulary.go, or the generated schema and table will not know it exists", kind)
		}
	}
}

// The reverse: the vocabulary must not document a kind nothing declares.
// That is how a list starts describing types that no longer exist.
func TestVocabularyDocumentsNothingThatIsNotDeclared(t *testing.T) {
	declared := make(map[string]bool)
	for _, kind := range declaredEventKinds(t) {
		declared[kind] = true
	}
	for _, doc := range orchestrator.EventVocabulary() {
		if !declared[string(doc.Kind)] {
			t.Errorf("EventVocabulary documents %q, which is not declared in %s", doc.Kind, vocabularySourceFile)
		}
	}
}

// A parse that finds nothing would make both tests above vacuous.
func TestTheSourceScanActuallyFindsEventKinds(t *testing.T) {
	if got := len(declaredEventKinds(t)); got < 10 {
		t.Fatalf("found %d event kinds in %s — the AST scan is not working, so the checks above prove nothing",
			got, vocabularySourceFile)
	}
}

// Every documented kind must say something useful. An empty summary is a
// row in the generated table that teaches a reader nothing.
func TestEveryDocumentedKindCarriesASummaryAndRoadmapItem(t *testing.T) {
	for _, doc := range orchestrator.EventVocabulary() {
		if strings.TrimSpace(doc.Summary) == "" {
			t.Errorf("event kind %q has no summary", doc.Kind)
		}
		if !strings.HasPrefix(doc.Roadmap, "L") && !strings.HasPrefix(doc.Roadmap, "M") {
			t.Errorf("event kind %q names roadmap item %q, want an L- or M-numbered item", doc.Kind, doc.Roadmap)
		}
	}
}

// The committed schema and table must match what the vocabulary generates,
// the same drift check typed state has.
func TestGeneratedTelemetryArtifactsAreCommitted(t *testing.T) {
	schema, err := orchestrator.GenerateEventSchema()
	if err != nil {
		t.Fatalf("GenerateEventSchema: %v", err)
	}
	generated := map[string][]byte{
		orchestrator.RunEventSchemaFile: schema,
		orchestrator.RunEventTableFile:  orchestrator.GenerateEventTable(),
	}
	for name, want := range generated {
		path := filepath.Join("..", "..", orchestrator.TelemetrySchemaDir, name)
		got, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v — run: go run ./cmd/gen-schemas", path, err)
		}
		if string(got) != string(want) {
			t.Errorf("%s has drifted from the vocabulary — run: go run ./cmd/gen-schemas", name)
		}
	}
}

// The schema must reject a kind nothing emits — otherwise `kind` is just a
// string and the vocabulary constrains nothing.
func TestGeneratedSchemaConstrainsKindToTheVocabulary(t *testing.T) {
	schema, err := orchestrator.GenerateEventSchema()
	if err != nil {
		t.Fatalf("GenerateEventSchema: %v", err)
	}
	text := string(schema)
	for _, kind := range orchestrator.EventKindStrings() {
		if !strings.Contains(text, strconv.Quote(kind)) {
			t.Errorf("generated schema does not list %q in the kind enum", kind)
		}
	}
	// A type specified elsewhere and still emitted by nothing. This used to
	// name policy.evaluated, which graduated in L2.16 when something finally
	// emitted it — which is the rule working: a type joins the vocabulary by
	// acquiring an emitter, never by being documented.
	if strings.Contains(text, strconv.Quote("contract.retry")) {
		t.Error("generated schema lists contract.retry, which nothing emits — the vocabulary holds only emitted types")
	}
}
