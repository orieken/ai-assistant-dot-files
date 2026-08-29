package tools

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
)

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

// reflectSchema derives a raw JSON Schema from a response struct.
// invopop/jsonschema stays an implementation detail of this package — the
// domain.Tool interface only ever sees stdlib json.RawMessage bytes. A marshal
// failure is a programmer error in the response struct, surfaced at
// registration via panic.
func reflectSchema(v any) json.RawMessage {
	body, err := json.Marshal(jsonschema.Reflect(v))
	if err != nil {
		panic("tools: reflected output schema failed to marshal: " + err.Error())
	}
	return body
}
