package server

import (
	"os"
	"path/filepath"

	"github.com/orieken/ai-assistant-dotfiles/mcp/internal/analyzers"
	"github.com/orieken/ai-assistant-dotfiles/mcp/internal/domain"
	"github.com/orieken/ai-assistant-dotfiles/mcp/internal/logging"
	"github.com/orieken/ai-assistant-dotfiles/mcp/internal/tools"
)

// buildFrameworkTools constructs and wires up the framework-generic tools.
func buildFrameworkTools(logger *logging.Logger) []domain.Tool {
	complexityAnalyzer := analyzers.NewComplexityAnalyzer()
	accessibilityAnalyzer := analyzers.NewAccessibilityAnalyzer()
	ubiquitousLanguageAnalyzer := analyzers.NewUbiquitousLanguageAnalyzer()
	dependencyBoundaryAnalyzer := analyzers.NewDependencyBoundaryAnalyzer()

	return []domain.Tool{
		tools.NewAnalyzeComplexityTool(logger, complexityAnalyzer),
		tools.NewCheckAccessibilityTool(logger, accessibilityAnalyzer),
		tools.NewCheckUbiquitousLanguageTool(logger, ubiquitousLanguageAnalyzer),
		tools.NewVerifyDependenciesTool(logger, dependencyBoundaryAnalyzer),
		tools.NewSearchKITool(logger, tools.NewKICorpusRetriever(frameworkCorpusPaths())),
		newSearchDocsTool(logger),
	}
}

func newSearchDocsTool(logger *logging.Logger) domain.Tool {
	dbPath, ok := docsFTSDBPath()
	if !ok {
		logger.Info("BM25 retriever not initialised — search_docs will return no-corpus diagnostic responses")
		return tools.NewSearchDocsTool(logger, nil, nil)
	}
	retriever, err := tools.NewBM25Retriever(dbPath)
	if err != nil {
		logger.Warn("BM25 retriever init failed — search_docs will return no-corpus diagnostic responses", "dbPath", dbPath, "error", err)
		return tools.NewSearchDocsTool(logger, nil, nil)
	}
	return tools.NewSearchDocsTool(logger, retriever, retriever)
}

func docsFTSDBPath() (string, bool) {
	if override := os.Getenv("DOCS_FTS_PATH"); override != "" {
		return override, true
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	if _, err := os.Stat(filepath.Join(cwd, ".claude")); err != nil {
		return "", false
	}
	return filepath.Join(cwd, ".claude", "rag", "docs-fts5.sqlite"), true
}

func frameworkCorpusPaths() []string {
	root := os.Getenv("AI_ASSISTANT_DOTFILES_PATH")
	if root == "" {
		return nil
	}
	return []string{
		filepath.Join(root, "shared", "knowledge"),
		filepath.Join(root, "docs", "adrs"),
	}
}
