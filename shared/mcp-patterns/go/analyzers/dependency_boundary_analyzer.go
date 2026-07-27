//go:build ignore
// +build ignore

// Code extracted from saturday-mcp (commit 0e5549125b7129e2b308df09d99e10d0b29a41bb).
// Reference implementation for framework-generic dependency boundary analyzer.
//
// When copying into downstream project:
// 1. Remove build tags above if compiling directly in your Go module.
// 2. Adjust package name and import paths (<YOUR_MODULE>/*) to match your codebase.

// Package analyzers - dependency_boundary_analyzer.go
//
// DependencyBoundaryAnalyzer is the concrete analyzer behind the
// verify_dependencies MCP tool (mcp-expand M1 Op 5). It walks a project
// tree, classifies each Go / TypeScript source file into a Clean
// Architecture layer based on its path segments, parses its import
// statements, and flags every import that crosses a layer boundary in
// the wrong direction (inner layer importing an outer layer).
package analyzers

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type DependencyBoundaryViolation struct {
	File       string `json:"file"`
	LineNumber int    `json:"lineNumber"`
	FromLayer  string `json:"fromLayer"`
	ToLayer    string `json:"toLayer"`
	ImportPath string `json:"importPath"`
}

type DependencyVerificationResult struct {
	Success         bool                          `json:"success"`
	ProjectPath     string                        `json:"projectPath"`
	ViolationsCount int                           `json:"violationsCount"`
	Violations      []DependencyBoundaryViolation `json:"violations,omitempty"`
	Summary         string                        `json:"summary"`
}

type DependencyBoundaryAnalyzer struct{}

func NewDependencyBoundaryAnalyzer() *DependencyBoundaryAnalyzer {
	return &DependencyBoundaryAnalyzer{}
}

const (
	layerUnknown    = "unknown"
	layerDomain     = "domain"
	layerUsecases   = "usecases"
	layerAdapters   = "adapters"
	layerFrameworks = "frameworks"
)

var layerRank = map[string]int{
	layerDomain:     0,
	layerUsecases:   1,
	layerAdapters:   2,
	layerFrameworks: 3,
}

var layerSegments = map[string]string{
	"domain":         layerDomain,
	"entities":       layerDomain,
	"entity":         layerDomain,
	"usecases":       layerUsecases,
	"use-cases":      layerUsecases,
	"use_cases":      layerUsecases,
	"usecase":        layerUsecases,
	"application":    layerUsecases,
	"adapters":       layerAdapters,
	"adapter":        layerAdapters,
	"interfaces":     layerAdapters,
	"infrastructure": layerAdapters,
	"frameworks":     layerFrameworks,
}

var scannedImportExtensions = map[string]struct{}{
	".go":  {},
	".ts":  {},
	".tsx": {},
}

var (
	goSingleImportRe    = regexp.MustCompile(`^\s*import\s+(?:[A-Za-z_.][\w]*\s+)?"([^"]+)"`)
	goImportBlockOpenRe = regexp.MustCompile(`^\s*import\s*\(\s*$`)
	goImportBlockLineRe = regexp.MustCompile(`^\s*(?:[A-Za-z_.][\w]*\s+)?"([^"]+)"\s*(?://.*)?$`)
	tsFromImportRe      = regexp.MustCompile(`from\s+['"]([^'"]+)['"]`)
	tsSideImportRe      = regexp.MustCompile(`^\s*import\s+['"]([^'"]+)['"]`)
)

type importRef struct {
	Path string
	Line int
}

func (a *DependencyBoundaryAnalyzer) Analyze(projectPath string) (*DependencyVerificationResult, error) {
	result := &DependencyVerificationResult{
		Success:     true,
		ProjectPath: projectPath,
		Violations:  []DependencyBoundaryViolation{},
	}

	files, err := a.collectSourceFiles(projectPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		a.scanFile(f, result)
	}

	result.ViolationsCount = len(result.Violations)
	result.Summary = summarizeDependencies(result.ViolationsCount)
	return result, nil
}

func summarizeDependencies(count int) string {
	if count == 0 {
		return "No dependency boundary violations found"
	}
	return "Dependency boundary violations found"
}

func classifyLayer(path string) string {
	parts := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == '\\'
	})
	for _, p := range parts {
		if l, ok := layerSegments[strings.ToLower(p)]; ok {
			return l
		}
	}
	return layerUnknown
}

func isBoundaryViolation(from, to string) bool {
	fromRank, ok1 := layerRank[from]
	toRank, ok2 := layerRank[to]
	if !ok1 || !ok2 {
		return false
	}
	return fromRank < toRank
}

func (a *DependencyBoundaryAnalyzer) collectSourceFiles(root string) ([]string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return []string{root}, nil
	}
	var files []string
	err = filepath.Walk(root, func(p string, entry os.FileInfo, walkErr error) error {
		if walkErr != nil || entry == nil {
			return nil
		}
		if entry.IsDir() {
			return SkipUninterestingDir(root, p, entry.Name())
		}
		if _, ok := scannedImportExtensions[strings.ToLower(filepath.Ext(p))]; ok {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (a *DependencyBoundaryAnalyzer) scanFile(file string, result *DependencyVerificationResult) {
	fromLayer := classifyLayer(file)
	if fromLayer == layerUnknown {
		return
	}
	imports, err := readImports(file)
	if err != nil {
		return
	}
	for _, imp := range imports {
		appendIfViolation(file, fromLayer, imp, result)
	}
}

func appendIfViolation(file, fromLayer string, imp importRef, result *DependencyVerificationResult) {
	toLayer := classifyLayer(imp.Path)
	if !isBoundaryViolation(fromLayer, toLayer) {
		return
	}
	result.Violations = append(result.Violations, DependencyBoundaryViolation{
		File:       file,
		LineNumber: imp.Line,
		FromLayer:  fromLayer,
		ToLayer:    toLayer,
		ImportPath: imp.Path,
	})
}

func readImports(file string) ([]importRef, error) {
	switch strings.ToLower(filepath.Ext(file)) {
	case ".go":
		return readGoImports(file)
	case ".ts", ".tsx":
		return readTSImports(file)
	}
	return nil, nil
}

func readGoImports(file string) ([]importRef, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var imports []importRef
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	lineNum := 0
	inBlock := false
	for scanner.Scan() {
		lineNum++
		inBlock = processGoImportLine(scanner.Text(), lineNum, inBlock, &imports)
	}
	return imports, nil
}

func processGoImportLine(line string, lineNum int, inBlock bool, imports *[]importRef) bool {
	if inBlock {
		return processGoImportBlockLine(line, lineNum, imports)
	}
	if goImportBlockOpenRe.MatchString(line) {
		return true
	}
	if m := goSingleImportRe.FindStringSubmatch(line); m != nil {
		*imports = append(*imports, importRef{Path: m[1], Line: lineNum})
	}
	return false
}

func processGoImportBlockLine(line string, lineNum int, imports *[]importRef) bool {
	if strings.TrimSpace(line) == ")" {
		return false
	}
	if m := goImportBlockLineRe.FindStringSubmatch(line); m != nil {
		*imports = append(*imports, importRef{Path: m[1], Line: lineNum})
	}
	return true
}

func readTSImports(file string) ([]importRef, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var imports []importRef
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		appendTSImportsFromLine(scanner.Text(), lineNum, &imports)
	}
	return imports, nil
}

func appendTSImportsFromLine(line string, lineNum int, imports *[]importRef) {
	if m := tsFromImportRe.FindStringSubmatch(line); m != nil {
		*imports = append(*imports, importRef{Path: m[1], Line: lineNum})
		return
	}
	if m := tsSideImportRe.FindStringSubmatch(line); m != nil {
		*imports = append(*imports, importRef{Path: m[1], Line: lineNum})
	}
}
