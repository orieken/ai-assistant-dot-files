package analyzers

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type FunctionComplexity struct {
	File         string `json:"file"`
	FunctionName string `json:"functionName"`
	LineNumber   int    `json:"lineNumber"`
	Complexity   int    `json:"complexity"`
	LineCount    int    `json:"lineCount"`
	Status       string `json:"status"`
}

type ComplexityAnalysisResult struct {
	Success         bool                 `json:"success"`
	ProjectPath     string               `json:"projectPath"`
	TotalFiles      int                  `json:"totalFiles"`
	TotalFunctions  int                  `json:"totalFunctions"`
	ViolationsCount int                  `json:"violationsCount"`
	Violations      []FunctionComplexity `json:"violations,omitempty"`
	Summary         string               `json:"summary"`
}

type ComplexityAnalyzer struct{}

func NewComplexityAnalyzer() *ComplexityAnalyzer { return &ComplexityAnalyzer{} }

func (a *ComplexityAnalyzer) Analyze(projectPath string, maxComplexity, maxLines int) (*ComplexityAnalysisResult, error) {
	if maxComplexity <= 0 {
		maxComplexity = 7
	}
	if maxLines <= 0 {
		maxLines = 30
	}

	result := &ComplexityAnalysisResult{
		Success:     true,
		ProjectPath: projectPath,
		Violations:  []FunctionComplexity{},
	}

	info, err := os.Stat(projectPath)
	if err != nil {
		return nil, err
	}

	var files []string
	if !info.IsDir() {
		files = append(files, projectPath)
	} else {
		err = filepath.Walk(projectPath, func(path string, fi os.FileInfo, err error) error {
			if err != nil || fi.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".ts") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	result.TotalFiles = len(files)
	for _, file := range files {
		if strings.HasSuffix(file, ".go") {
			a.analyzeGoFile(file, maxComplexity, maxLines, result)
		} else {
			a.analyzeGenericFile(file, maxComplexity, maxLines, result)
		}
	}

	if result.ViolationsCount > 0 {
		result.Summary = "Complexity threshold violations found"
	} else {
		result.Summary = "All functions pass complexity and LOC checks"
	}
	return result, nil
}

func (a *ComplexityAnalyzer) analyzeGoFile(file string, maxComplexity, maxLines int, result *ComplexityAnalysisResult) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		a.analyzeGenericFile(file, maxComplexity, maxLines, result)
		return
	}

	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		result.TotalFunctions++

		startPos := fset.Position(fn.Pos())
		endPos := fset.Position(fn.End())
		lineCount := endPos.Line - startPos.Line + 1

		complexity := 1
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
				_ = v
				complexity++
			case *ast.BinaryExpr:
				if v.Op == token.LAND || v.Op == token.LOR {
					complexity++
				}
			}
			return true
		})

		if complexity > maxComplexity || lineCount > maxLines {
			result.ViolationsCount++
			result.Violations = append(result.Violations, FunctionComplexity{
				File:         file,
				FunctionName: fn.Name.Name,
				LineNumber:   startPos.Line,
				Complexity:   complexity,
				LineCount:    lineCount,
				Status:       "VIOLATION",
			})
		}
	}
}

func (a *ComplexityAnalyzer) analyzeGenericFile(file string, maxComplexity, maxLines int, result *ComplexityAnalysisResult) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lines := 0
	for scanner.Scan() {
		lines++
	}
	result.TotalFunctions++
	if lines > maxLines {
		result.ViolationsCount++
		result.Violations = append(result.Violations, FunctionComplexity{
			File:         file,
			FunctionName: filepath.Base(file),
			LineNumber:   1,
			Complexity:   1,
			LineCount:    lines,
			Status:       "VIOLATION",
		})
	}
}
