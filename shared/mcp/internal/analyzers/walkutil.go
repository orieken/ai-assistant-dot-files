package analyzers

import (
	"path/filepath"
	"strings"
)

// SkippedDirNames lists directory basenames the walk should never descend into.
var SkippedDirNames = map[string]struct{}{
	"node_modules": {},
	".git":         {},
	"dist":         {},
	"build":        {},
	"vendor":       {},
	".venv":        {},
}

// SkipUninterestingDir returns filepath.SkipDir for directories the walk should skip:
// any name in SkippedDirNames, plus hidden dot-directories below the root.
func SkipUninterestingDir(root, current, name string) error {
	if _, skip := SkippedDirNames[name]; skip {
		return filepath.SkipDir
	}
	if strings.HasPrefix(name, ".") && current != root {
		return filepath.SkipDir
	}
	return nil
}
