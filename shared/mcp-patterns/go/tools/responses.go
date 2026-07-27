//go:build ignore
// +build ignore

// Code extracted from saturday-mcp (commit 0e5549125b7129e2b308df09d99e10d0b29a41bb).
// Reference implementation for framework-generic M1 tool responses and schema reflection.
//
// When copying into downstream project:
// 1. Remove build tags above if compiling directly in your Go module.
// 2. Adjust package name and import paths (<YOUR_MODULE>/*) to match your codebase.

package tools

import "github.com/invopop/jsonschema"

// FunctionComplexity captures complexity metrics for a single function.
type FunctionComplexity struct {
	File         string `json:"file"`
	FunctionName string `json:"functionName"`
	LineNumber   int    `json:"lineNumber"`
	Complexity   int    `json:"complexity"`
	LineCount    int    `json:"lineCount"`
	Status       string `json:"status"`
}

// ComplexityAnalysisResult is the response from analyze_complexity.
type ComplexityAnalysisResult struct {
	Success         bool                 `json:"success"`
	ProjectPath     string               `json:"projectPath"`
	TotalFiles      int                  `json:"totalFiles"`
	TotalFunctions  int                  `json:"totalFunctions"`
	ViolationsCount int                  `json:"violationsCount"`
	Violations      []FunctionComplexity `json:"violations,omitempty"`
	Summary         string               `json:"summary"`
}

// AccessibilityViolation details a single a11y defect.
type AccessibilityViolation struct {
	File        string `json:"file"`
	LineNumber  int    `json:"lineNumber"`
	Element     string `json:"element"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

// AccessibilityReportResult is the response from check_accessibility.
type AccessibilityReportResult struct {
	Success         bool                     `json:"success"`
	Path            string                   `json:"path"`
	TotalFiles      int                      `json:"totalFiles"`
	ViolationsCount int                      `json:"violationsCount"`
	Violations      []AccessibilityViolation `json:"violations,omitempty"`
	Summary         string                   `json:"summary"`
}

// LanguageViolation represents a domain term violation.
type LanguageViolation struct {
	File        string `json:"file"`
	LineNumber  int    `json:"lineNumber"`
	InvalidTerm string `json:"invalidTerm"`
	Suggested   string `json:"suggested"`
}

// UbiquitousLanguageResult is the response from check_ubiquitous_language.
type UbiquitousLanguageResult struct {
	Success         bool                `json:"success"`
	ProjectPath     string              `json:"projectPath"`
	ViolationsCount int                 `json:"violationsCount"`
	Violations      []LanguageViolation `json:"violations,omitempty"`
	Summary         string              `json:"summary"`
}

// DependencyBoundaryViolation represents a Clean Architecture layer violation.
type DependencyBoundaryViolation struct {
	File       string `json:"file"`
	LineNumber int    `json:"lineNumber"`
	FromLayer  string `json:"fromLayer"`
	ToLayer    string `json:"toLayer"`
	ImportPath string `json:"importPath"`
}

// DependencyVerificationResult is the response from verify_dependencies.
type DependencyVerificationResult struct {
	Success         bool                          `json:"success"`
	ProjectPath     string                        `json:"projectPath"`
	ViolationsCount int                           `json:"violationsCount"`
	Violations      []DependencyBoundaryViolation `json:"violations,omitempty"`
	Summary         string                        `json:"summary"`
}

// KIMatch represents a Knowledge Item search hit.
type KIMatch struct {
	Title     string   `json:"title"`
	Path      string   `json:"path"`
	Summary   string   `json:"summary"`
	Tags      []string `json:"tags,omitempty"`
	Relevance float64  `json:"relevance"`
}

// KISearchResult is the response from search_ki.
type KISearchResult struct {
	Success   bool      `json:"success"`
	Query     string    `json:"query"`
	TotalHits int       `json:"totalHits"`
	Matches   []KIMatch `json:"matches,omitempty"`
}

// DocMatch represents a single BM25-ranked docs-corpus hit.
type DocMatch struct {
	Title     string  `json:"title"`
	Path      string  `json:"path"`
	Summary   string  `json:"summary"`
	Relevance float64 `json:"relevance"`
}

// DocSearchResult is the response from search_docs.
type DocSearchResult struct {
	Success   bool       `json:"success"`
	Query     string     `json:"query"`
	TotalHits int        `json:"totalHits"`
	Matches   []DocMatch `json:"matches,omitempty"`
}

// reflectSchema wraps jsonschema.Reflect so every OutputSchema() reads the same one-liner.
func reflectSchema(v interface{}) *jsonschema.Schema {
	return jsonschema.Reflect(v)
}
