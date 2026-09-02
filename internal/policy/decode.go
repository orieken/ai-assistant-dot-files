package policy

// Walking a YAML condition node into typed checks.
//
// This is where the absence of an expression language costs something and
// buys something. It costs this file: a hand-written walk over a small
// grammar. It buys a policy author an error that names the field they got
// wrong, and a reader of the code a complete list of what a policy can
// possibly do.

import (
	"fmt"
	"strconv"

	"go.yaml.in/yaml/v4"
)

// The three structural keys, which are not fields.
const (
	keyAny = "any"
	keyNot = "not"
)

func decodeCondition(node *yaml.Node) (Condition, error) {
	if node == nil || node.Kind == 0 {
		return Condition{}, nil
	}
	if node.Kind != yaml.MappingNode {
		return Condition{}, fmt.Errorf("condition must be a mapping, got %s", nodeKind(node))
	}
	var condition Condition
	seen := make(map[string]bool, len(node.Content)/2)
	for index := 0; index+1 < len(node.Content); index += 2 {
		key := node.Content[index].Value
		if err := checkNotRepeated(seen, key); err != nil {
			return Condition{}, err
		}
		if err := addEntry(&condition, key, node.Content[index+1]); err != nil {
			return Condition{}, err
		}
	}
	return condition, nil
}

// checkNotRepeated rejects a field appearing twice in one condition
// mapping.
//
// The struct decoder catches duplicates at the document's top level, but
// the condition arrives as a raw node so it can be walked into typed
// checks — and a raw node keeps every duplicate key silently. That is
// precisely where the bug was: the shipped auto-approve-refactor policy
// listed `filePaths` twice, and a lenient reader would take the second and
// discard the first, so the security-path exclusion it advertised would
// never have run. In an authorization document a silently discarded check
// is a check that does not exist.
func checkNotRepeated(seen map[string]bool, key string) error {
	if seen[key] {
		return fmt.Errorf("condition field %q appears more than once — a repeated key silently discards "+
			"the earlier check; combine alternatives under `not: {any: [...]}` instead", key)
	}
	seen[key] = true
	return nil
}

// addEntry dispatches one key of a condition mapping.
func addEntry(condition *Condition, key string, value *yaml.Node) error {
	switch key {
	case keyAny:
		return addAny(condition, value)
	case keyNot:
		return addNot(condition, value)
	default:
		return addChecks(condition, Field(key), value)
	}
}

func addAny(condition *Condition, value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return fmt.Errorf("`any` must be a list of conditions, got %s", nodeKind(value))
	}
	for _, item := range value.Content {
		branch, err := decodeCondition(item)
		if err != nil {
			return err
		}
		condition.Any = append(condition.Any, branch)
	}
	return nil
}

func addNot(condition *Condition, value *yaml.Node) error {
	inner, err := decodeCondition(value)
	if err != nil {
		return err
	}
	condition.Not = &inner
	return nil
}

// addChecks handles both shapes a field takes: a bare scalar
// (`testsPass: true`) and a mapping of operators (`diffLines: {lessThan: 200}`).
func addChecks(condition *Condition, field Field, value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		check, err := scalarCheck(field, value)
		if err != nil {
			return err
		}
		condition.Checks = append(condition.Checks, check)
		return nil
	}
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("field %q must be a value or a mapping of operators, got %s", field, nodeKind(value))
	}
	return addOperatorChecks(condition, field, value)
}

func addOperatorChecks(condition *Condition, field Field, value *yaml.Node) error {
	for index := 0; index+1 < len(value.Content); index += 2 {
		check, err := operatorCheck(field, Operator(value.Content[index].Value), value.Content[index+1])
		if err != nil {
			return err
		}
		condition.Checks = append(condition.Checks, check)
	}
	return nil
}

// scalarCheck reads `field: true`, whose operator is implicit.
func scalarCheck(field Field, value *yaml.Node) (Check, error) {
	boolean, err := strconv.ParseBool(value.Value)
	if err != nil {
		return Check{}, fmt.Errorf("field %q used as a bare value must be true or false, got %q", field, value.Value)
	}
	if err := validateCheck(field, OpIsTrue, kindBool); err != nil {
		return Check{}, err
	}
	return Check{Field: field, Operator: OpIsTrue, Bool: boolean}, nil
}

func operatorCheck(field Field, operator Operator, value *yaml.Node) (Check, error) {
	kind, err := kindOf(value)
	if err != nil {
		return Check{}, fmt.Errorf("field %q, operator %q: %w", field, operator, err)
	}
	if err := validateCheck(field, operator, kind); err != nil {
		return Check{}, err
	}
	return buildCheck(field, operator, kind, value)
}

func buildCheck(field Field, operator Operator, kind valueKind, value *yaml.Node) (Check, error) {
	check := Check{Field: field, Operator: operator}
	switch kind {
	case kindNumber:
		number, err := strconv.ParseFloat(value.Value, 64)
		if err != nil {
			return Check{}, fmt.Errorf("field %q: %q is not a number", field, value.Value)
		}
		check.Number = number
	case kindBool:
		check.Bool = value.Value == "true"
	case kindText:
		check.Text = value.Value
	}
	return check, nil
}

// kindOf classifies a scalar. YAML tags are authoritative here rather than
// guessing from the text, so the quoted string "200" is not silently a
// number.
func kindOf(node *yaml.Node) (valueKind, error) {
	if node.Kind != yaml.ScalarNode {
		return "", fmt.Errorf("expected a value, got %s", nodeKind(node))
	}
	switch node.Tag {
	case "!!int", "!!float":
		return kindNumber, nil
	case "!!bool":
		return kindBool, nil
	case "!!str":
		return kindText, nil
	default:
		return "", fmt.Errorf("unsupported value %q", node.Value)
	}
}

func nodeKind(node *yaml.Node) string {
	names := map[yaml.Kind]string{
		yaml.DocumentNode: "a document", yaml.SequenceNode: "a list",
		yaml.MappingNode: "a mapping", yaml.ScalarNode: "a value", yaml.AliasNode: "an alias",
	}
	if name, ok := names[node.Kind]; ok {
		return name
	}
	return "nothing"
}
