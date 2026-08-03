package tools

import "github.com/invopop/jsonschema"

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
