package policy

// `**` matching, which filepath.Match does not support.
//
// The schema's own examples depend on it — `**/security/**`,
// `{tests,test,spec}/**` — so a matcher without it would silently fail
// every path exclusion in the shipped policies. Silently, because a
// non-matching exclusion reads exactly like a clean run.

import "strings"

// matchesDoubleStar splits the pattern on `**` and walks the path,
// requiring each literal segment to appear in order.
func matchesDoubleStar(pattern, path string) bool {
	if !strings.Contains(pattern, "**") {
		return false
	}
	parts := strings.Split(pattern, "**")
	return consumeParts(parts, path)
}

func consumeParts(parts []string, path string) bool {
	remainder := path
	for index, part := range parts {
		trimmed := strings.Trim(part, "/")
		if trimmed == "" {
			continue
		}
		next, ok := consumeOne(trimmed, remainder, index == 0)
		if !ok {
			return false
		}
		remainder = next
	}
	return true
}

// consumeOne finds the next literal segment. The first part is anchored to
// the start of the path — `docs/**` must not match `src/docs/x` — while
// later parts may appear anywhere after what has already been consumed.
func consumeOne(part, remainder string, anchored bool) (string, bool) {
	if anchored {
		if !strings.HasPrefix(remainder, part) {
			return "", false
		}
		return remainder[len(part):], true
	}
	index := strings.Index(remainder, part)
	if index < 0 {
		return "", false
	}
	return remainder[index+len(part):], true
}
