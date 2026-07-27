//go:build ignore
// +build ignore

// Code extracted from saturday-mcp (commit 0e5549125b7129e2b308df09d99e10d0b29a41bb).
// Reference implementation for framework-generic accessibility analyzer.
//
// When copying into downstream project:
// 1. Remove build tags above if compiling directly in your Go module.
// 2. Adjust package name and import paths (<YOUR_MODULE>/*) to match your codebase.

// Package analyzers - accessibility_analyzer.go
//
// AccessibilityAnalyzer is the concrete analyzer behind the
// check_accessibility MCP tool (mcp-expand M1 Op 3). It walks a project
// path (or scans a single file) matching common UI template extensions
// and reports semantic-HTML / ARIA violations line by line.
package analyzers

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// AccessibilityViolation captures a single a11y defect the analyzer found.
type AccessibilityViolation struct {
	File        string `json:"file"`
	LineNumber  int    `json:"lineNumber"`
	Element     string `json:"element"`
	Rule        string `json:"rule"`
	Description string `json:"description"`
}

// AccessibilityReportResult is the full analyzer response.
type AccessibilityReportResult struct {
	Success         bool                     `json:"success"`
	Path            string                   `json:"path"`
	TotalFiles      int                      `json:"totalFiles"`
	ViolationsCount int                      `json:"violationsCount"`
	Violations      []AccessibilityViolation `json:"violations,omitempty"`
	Summary         string                   `json:"summary"`
}

// AccessibilityAnalyzer scans UI template files for semantic-HTML and ARIA violations.
type AccessibilityAnalyzer struct{}

// NewAccessibilityAnalyzer constructs an analyzer with the default rule set.
func NewAccessibilityAnalyzer() *AccessibilityAnalyzer {
	return &AccessibilityAnalyzer{}
}

var scannedExtensions = map[string]struct{}{
	".html":   {},
	".htm":    {},
	".vue":    {},
	".jsx":    {},
	".tsx":    {},
	".svelte": {},
}

const (
	ruleClickableNonInteractive = "clickable-non-interactive"
	ruleImgAlt                  = "img-alt"
	ruleAnchorName              = "anchor-name"
	ruleInputLabel              = "input-label"
	ruleHTMLLang                = "html-lang"
)

var (
	clickableNonInteractiveRe = regexp.MustCompile(`(?i)<(div|span)\b[^>]*\bon[a-z]+\s*=`)
	imgTagRe                  = regexp.MustCompile(`(?i)<img\b[^>]*>`)
	imgHasAltRe               = regexp.MustCompile(`(?i)\balt\s*=`)
	anchorOpenRe              = regexp.MustCompile(`(?i)<a\b[^>]*>`)
	ariaLabelRe               = regexp.MustCompile(`(?i)\baria-label\s*=\s*["'][^"']+["']`)
	inputTagRe                = regexp.MustCompile(`(?i)<input\b[^>]*>`)
	inputHasIDRe              = regexp.MustCompile(`(?i)\bid\s*=\s*["']([^"']+)["']`)
	inputTypeHiddenRe         = regexp.MustCompile(`(?i)\btype\s*=\s*["']hidden["']`)
	labelForRe                = regexp.MustCompile(`(?i)<label\b[^>]*\bfor\s*=\s*["']([^"']+)["']`)
	htmlTagRe                 = regexp.MustCompile(`(?i)<html\b[^>]*>`)
	htmlHasLangRe             = regexp.MustCompile(`(?i)\blang\s*=`)
)

func (a *AccessibilityAnalyzer) Analyze(path string) (*AccessibilityReportResult, error) {
	result := &AccessibilityReportResult{
		Success:    true,
		Path:       path,
		Violations: []AccessibilityViolation{},
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	files, err := a.collectFiles(path, info.IsDir())
	if err != nil {
		return nil, err
	}
	result.TotalFiles = len(files)

	for _, f := range files {
		a.analyzeFile(f, result)
	}

	result.ViolationsCount = len(result.Violations)
	result.Summary = summarize(result.ViolationsCount)
	return result, nil
}

func summarize(count int) string {
	if count == 0 {
		return "No accessibility violations found"
	}
	return "Accessibility violations found"
}

func (a *AccessibilityAnalyzer) collectFiles(path string, isDir bool) ([]string, error) {
	if !isDir {
		return []string{path}, nil
	}
	var files []string
	err := filepath.Walk(path, func(p string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info == nil {
			return nil
		}
		if info.IsDir() {
			return SkipUninterestingDir(path, p, info.Name())
		}
		if _, ok := scannedExtensions[strings.ToLower(filepath.Ext(p))]; ok {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (a *AccessibilityAnalyzer) analyzeFile(file string, result *AccessibilityReportResult) {
	labels := collectLabelIDs(file)
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
		runLineRules(file, lineNum, scanner.Text(), labels, result)
	}
}

func runLineRules(file string, lineNum int, line string, labels map[string]struct{}, result *AccessibilityReportResult) {
	checkHTMLLang(file, lineNum, line, result)
	checkClickableNonInteractive(file, lineNum, line, result)
	checkImgAlt(file, lineNum, line, result)
	checkAnchorName(file, lineNum, line, result)
	checkInputLabel(file, lineNum, line, labels, result)
}

func checkHTMLLang(file string, lineNum int, line string, result *AccessibilityReportResult) {
	if !htmlTagRe.MatchString(line) || htmlHasLangRe.MatchString(line) {
		return
	}
	result.Violations = append(result.Violations, AccessibilityViolation{
		File:        file,
		LineNumber:  lineNum,
		Element:     "html",
		Rule:        ruleHTMLLang,
		Description: "<html> element is missing a lang attribute",
	})
}

func checkClickableNonInteractive(file string, lineNum int, line string, result *AccessibilityReportResult) {
	m := clickableNonInteractiveRe.FindStringSubmatch(line)
	if m == nil {
		return
	}
	el := strings.ToLower(m[1])
	result.Violations = append(result.Violations, AccessibilityViolation{
		File:        file,
		LineNumber:  lineNum,
		Element:     el,
		Rule:        ruleClickableNonInteractive,
		Description: "Non-interactive element (<" + el + ">) has an event handler; use <button> or <a> instead",
	})
}

func checkImgAlt(file string, lineNum int, line string, result *AccessibilityReportResult) {
	tag := imgTagRe.FindString(line)
	if tag == "" || imgHasAltRe.MatchString(tag) {
		return
	}
	result.Violations = append(result.Violations, AccessibilityViolation{
		File:        file,
		LineNumber:  lineNum,
		Element:     "img",
		Rule:        ruleImgAlt,
		Description: "<img> element is missing an alt attribute",
	})
}

func checkAnchorName(file string, lineNum int, line string, result *AccessibilityReportResult) {
	openTag := anchorOpenRe.FindString(line)
	if openTag == "" || ariaLabelRe.MatchString(openTag) {
		return
	}
	if hasVisibleAnchorText(line, openTag) {
		return
	}
	result.Violations = append(result.Violations, AccessibilityViolation{
		File:        file,
		LineNumber:  lineNum,
		Element:     "a",
		Rule:        ruleAnchorName,
		Description: "<a> element has no visible text and no aria-label",
	})
}

func hasVisibleAnchorText(line, openTag string) bool {
	rest := line[strings.Index(line, openTag)+len(openTag):]
	closeIdx := strings.Index(strings.ToLower(rest), "</a>")
	if closeIdx < 0 {
		return strings.TrimSpace(rest) != ""
	}
	return strings.TrimSpace(rest[:closeIdx]) != ""
}

func checkInputLabel(file string, lineNum int, line string, labels map[string]struct{}, result *AccessibilityReportResult) {
	tag := inputTagRe.FindString(line)
	if tag == "" || inputTypeHiddenRe.MatchString(tag) {
		return
	}
	if ariaLabelRe.MatchString(tag) || hasMatchingLabel(tag, labels) {
		return
	}
	result.Violations = append(result.Violations, AccessibilityViolation{
		File:        file,
		LineNumber:  lineNum,
		Element:     "input",
		Rule:        ruleInputLabel,
		Description: "<input> element has no associated <label for=...> and no aria-label",
	})
}

func hasMatchingLabel(tag string, labels map[string]struct{}) bool {
	idMatch := inputHasIDRe.FindStringSubmatch(tag)
	if idMatch == nil {
		return false
	}
	_, ok := labels[idMatch[1]]
	return ok
}

func collectLabelIDs(file string) map[string]struct{} {
	ids := map[string]struct{}{}
	f, err := os.Open(file)
	if err != nil {
		return ids
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		for _, m := range labelForRe.FindAllStringSubmatch(scanner.Text(), -1) {
			ids[m[1]] = struct{}{}
		}
	}
	return ids
}
