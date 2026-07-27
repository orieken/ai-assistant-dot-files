//go:build ignore
// +build ignore

// Code extracted from saturday-mcp (commit 0e5549125b7129e2b308df09d99e10d0b29a41bb).
// Reference implementation for framework-generic ubiquitous language analyzer.
//
// When copying into downstream project:
// 1. Remove build tags above if compiling directly in your Go module.
// 2. Adjust package name and import paths (<YOUR_MODULE>/*) to match your codebase.

// Package analyzers - ubiquitous_language_analyzer.go
//
// UbiquitousLanguageAnalyzer is the concrete analyzer behind the
// check_ubiquitous_language MCP tool (mcp-expand M1 Op 4). It parses a
// DOMAIN_DICTIONARY.md, extracts canonical terms plus their forbidden
// synonyms, walks a project tree, and reports every source-file line
// that mentions any synonym alongside the canonical replacement.
package analyzers

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type LanguageViolation struct {
	File        string `json:"file"`
	LineNumber  int    `json:"lineNumber"`
	InvalidTerm string `json:"invalidTerm"`
	Suggested   string `json:"suggested"`
}

type UbiquitousLanguageResult struct {
	Success         bool                `json:"success"`
	ProjectPath     string              `json:"projectPath"`
	ViolationsCount int                 `json:"violationsCount"`
	Violations      []LanguageViolation `json:"violations,omitempty"`
	Summary         string              `json:"summary"`
}

type UbiquitousLanguageAnalyzer struct{}

func NewUbiquitousLanguageAnalyzer() *UbiquitousLanguageAnalyzer {
	return &UbiquitousLanguageAnalyzer{}
}

var scannedSourceExtensions = map[string]struct{}{
	".go":   {},
	".ts":   {},
	".tsx":  {},
	".js":   {},
	".jsx":  {},
	".py":   {},
	".java": {},
	".cs":   {},
}

var (
	backtickTermRe   = regexp.MustCompile("`([^`]+)`")
	sectionHeaderRe  = regexp.MustCompile(`^##\s+(.+?)\s*$`)
	numberedHeaderRe = regexp.MustCompile(`^\d+\.\s`)
	synonymLineRe    = regexp.MustCompile(`(?i)^\s*\*\*(?:synonyms to avoid|synonyms|avoid|not)\*\*\s*:\s*(.+)$`)
	tableRowRe       = regexp.MustCompile(`^\|\s*\*\*([^*]+)\*\*\s*\|(.+)$`)
)

type termMatcher struct {
	re        *regexp.Regexp
	synonym   string
	canonical string
}

func (a *UbiquitousLanguageAnalyzer) Analyze(projectPath, dictionaryPath string) (*UbiquitousLanguageResult, error) {
	result := &UbiquitousLanguageResult{
		Success:     true,
		ProjectPath: projectPath,
		Violations:  []LanguageViolation{},
	}

	synonyms, err := loadDictionary(dictionaryPath)
	if err != nil {
		return nil, err
	}
	if len(synonyms) == 0 {
		result.Summary = summarizeLanguage(0)
		return result, nil
	}

	files, err := a.collectSourceFiles(projectPath)
	if err != nil {
		return nil, err
	}

	matchers := compileMatchers(synonyms)
	for _, f := range files {
		a.scanFile(f, matchers, result)
	}

	result.ViolationsCount = len(result.Violations)
	result.Summary = summarizeLanguage(result.ViolationsCount)
	return result, nil
}

func summarizeLanguage(count int) string {
	if count == 0 {
		return "No ubiquitous language violations found"
	}
	return "Ubiquitous language violations found"
}

func loadDictionary(dictionaryPath string) (map[string]string, error) {
	f, err := os.Open(dictionaryPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	synonyms := map[string]string{}
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var currentTerm string
	for scanner.Scan() {
		line := scanner.Text()
		currentTerm = updateTermFromHeader(line, currentTerm)
		parseTableRow(line, synonyms)
		parseSynonymLine(line, currentTerm, synonyms)
	}
	return synonyms, nil
}

func updateTermFromHeader(line, current string) string {
	m := sectionHeaderRe.FindStringSubmatch(line)
	if m == nil {
		return current
	}
	if numberedHeaderRe.MatchString(m[1]) {
		return ""
	}
	return m[1]
}

func parseTableRow(line string, synonyms map[string]string) {
	m := tableRowRe.FindStringSubmatch(line)
	if m == nil {
		return
	}
	canonical := strings.TrimSpace(m[1])
	parts := splitTableCells(m[2])
	if len(parts) < 2 {
		return
	}
	addSynonyms(parts[len(parts)-1], canonical, synonyms)
}

func splitTableCells(raw string) []string {
	parts := strings.Split(raw, "|")
	for len(parts) > 0 && strings.TrimSpace(parts[len(parts)-1]) == "" {
		parts = parts[:len(parts)-1]
	}
	return parts
}

func parseSynonymLine(line, currentTerm string, synonyms map[string]string) {
	if currentTerm == "" {
		return
	}
	m := synonymLineRe.FindStringSubmatch(line)
	if m == nil {
		return
	}
	addSynonyms(m[1], currentTerm, synonyms)
}

func addSynonyms(raw, canonical string, synonyms map[string]string) {
	matches := backtickTermRe.FindAllStringSubmatch(raw, -1)
	if len(matches) > 0 {
		for _, m := range matches {
			recordSynonym(m[1], canonical, synonyms)
		}
		return
	}
	for _, s := range strings.Split(raw, ",") {
		recordSynonym(s, canonical, synonyms)
	}
}

func recordSynonym(term, canonical string, synonyms map[string]string) {
	term = strings.TrimSpace(term)
	term = strings.Trim(term, "`*_ ")
	if idx := strings.Index(term, "("); idx > 0 {
		term = strings.TrimSpace(term[:idx])
	}
	if term == "" {
		return
	}
	synonyms[strings.ToLower(term)] = canonical
}

func compileMatchers(synonyms map[string]string) []termMatcher {
	matchers := make([]termMatcher, 0, len(synonyms))
	for synonym, canonical := range synonyms {
		pattern := `(?i)\b` + regexp.QuoteMeta(synonym) + `\b`
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		matchers = append(matchers, termMatcher{re: re, synonym: synonym, canonical: canonical})
	}
	return matchers
}

func (a *UbiquitousLanguageAnalyzer) collectSourceFiles(root string) ([]string, error) {
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
		if _, ok := scannedSourceExtensions[strings.ToLower(filepath.Ext(p))]; ok {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (a *UbiquitousLanguageAnalyzer) scanFile(file string, matchers []termMatcher, result *UbiquitousLanguageResult) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		appendLineViolations(file, lineNum, scanner.Text(), matchers, result)
	}
}

func appendLineViolations(file string, lineNum int, line string, matchers []termMatcher, result *UbiquitousLanguageResult) {
	for _, m := range matchers {
		if m.re.MatchString(line) {
			result.Violations = append(result.Violations, LanguageViolation{
				File:        file,
				LineNumber:  lineNum,
				InvalidTerm: m.synonym,
				Suggested:   m.canonical,
			})
		}
	}
}
