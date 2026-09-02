package policy

// The condition vocabulary, typed. Every check a policy may express is a
// member of this file; there is no escape hatch, which is the point.

import (
	"fmt"
	"sort"
)

// Field names a fact a condition can test.
type Field string

// The fields policy-schema.md declares.
const (
	FieldDiffLines            Field = "diffLines"
	FieldDiffType             Field = "diffType"
	FieldTestsPass            Field = "testsPass"
	FieldDryRunPass           Field = "dryRunPass"
	FieldFilePaths            Field = "filePaths"
	FieldReviewVerdict        Field = "codeReviewer.verdict"
	FieldReviewBehaviorChange Field = "codeReviewer.behaviorChange"
	FieldSecurityCriticals    Field = "securityReviewer.criticals"
	FieldFitnessAllPass       Field = "fitnessFunction.allPass"
)

// Operator names how a field is compared.
type Operator string

// The operators policy-schema.md declares.
const (
	OpEquals    Operator = "equals"
	OpLessThan  Operator = "lessThan"
	OpAllMatch  Operator = "allMatch"
	OpNoneMatch Operator = "noneMatch"
	OpAnyMatch  Operator = "anyMatch"
	// OpIsTrue is the implicit operator of a bare boolean check such as
	// `testsPass: true`.
	OpIsTrue Operator = "isTrue"
)

// valueKind is what an operator's operand must be.
type valueKind string

const (
	kindNumber valueKind = "number"
	kindText   valueKind = "string"
	kindBool   valueKind = "bool"
)

// fieldRule declares what a field accepts. A table rather than a switch,
// for the same reason the state package's factories are one: adding a check
// is one row, and the error messages derive from the same data the parser
// enforces.
type fieldRule struct {
	operators map[Operator]valueKind
}

func fieldRules() map[Field]fieldRule {
	globs := map[Operator]valueKind{OpAllMatch: kindText, OpNoneMatch: kindText, OpAnyMatch: kindText}
	return map[Field]fieldRule{
		FieldDiffLines:            {operators: map[Operator]valueKind{OpLessThan: kindNumber, OpEquals: kindNumber}},
		FieldDiffType:             {operators: map[Operator]valueKind{OpEquals: kindText}},
		FieldTestsPass:            {operators: map[Operator]valueKind{OpIsTrue: kindBool, OpEquals: kindBool}},
		FieldDryRunPass:           {operators: map[Operator]valueKind{OpIsTrue: kindBool, OpEquals: kindBool}},
		FieldFilePaths:            {operators: globs},
		FieldReviewVerdict:        {operators: map[Operator]valueKind{OpEquals: kindText}},
		FieldReviewBehaviorChange: {operators: map[Operator]valueKind{OpIsTrue: kindBool, OpEquals: kindBool}},
		FieldSecurityCriticals:    {operators: map[Operator]valueKind{OpEquals: kindNumber, OpLessThan: kindNumber}},
		FieldFitnessAllPass:       {operators: map[Operator]valueKind{OpIsTrue: kindBool, OpEquals: kindBool}},
	}
}

// KnownFields returns every field a policy may test, sorted, for error
// messages that tell the author what they could have written.
func KnownFields() []Field {
	rules := fieldRules()
	fields := make([]Field, 0, len(rules))
	for field := range rules {
		fields = append(fields, field)
	}
	sort.Slice(fields, func(i, j int) bool { return fields[i] < fields[j] })
	return fields
}

// Check is one typed comparison. Exactly one of the value fields is set,
// determined by the operator.
type Check struct {
	Field    Field
	Operator Operator
	Number   float64
	Text     string
	Bool     bool
}

// Condition is a tree: checks combined with implicit AND, an optional Any
// (OR), and an optional Not.
type Condition struct {
	Checks []Check
	Any    []Condition
	Not    *Condition
}

// IsEmpty reports a condition that tests nothing. A policy with an empty
// condition would fire on every gate it matches, which is never what
// someone meant to write.
func (c Condition) IsEmpty() bool {
	return len(c.Checks) == 0 && len(c.Any) == 0 && c.Not == nil
}

// validateCheck confirms a field accepts an operator and that the operand
// is the right kind.
func validateCheck(field Field, operator Operator, kind valueKind) error {
	rule, known := fieldRules()[field]
	if !known {
		return fmt.Errorf("unknown condition field %q — valid fields are %v", field, KnownFields())
	}
	want, accepted := rule.operators[operator]
	if !accepted {
		return fmt.Errorf("field %q does not accept operator %q", field, operator)
	}
	if want != kind {
		return fmt.Errorf("field %q with %q expects a %s, got a %s", field, operator, want, kind)
	}
	return nil
}
