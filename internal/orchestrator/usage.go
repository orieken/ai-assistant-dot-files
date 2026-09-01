package orchestrator

// Usage is what one model invocation consumed (roadmap L3.8). Every number
// here is REPORTED by the provider, never computed here — the cost field
// especially. A pricing table in this repository would be wrong within a
// quarter and would produce a confident figure nobody had been charged,
// which is the same class of defect as a duration recalled by a model.
//
// A nil *Usage means the provider reported nothing, which is different from
// a provider reporting zeros. The mock provider reports nothing.
type Usage struct {
	// Model is what the provider says actually served the request, which is
	// not necessarily what was asked for.
	Model string `json:"model,omitempty"`
	// InputTokens and OutputTokens are the billable token counts.
	InputTokens  int64 `json:"inputTokens,omitempty"`
	OutputTokens int64 `json:"outputTokens,omitempty"`
	// CacheReadTokens and CacheCreationTokens are prompt-cache traffic,
	// kept separate because they price differently from fresh input.
	CacheReadTokens     int64 `json:"cacheReadTokens,omitempty"`
	CacheCreationTokens int64 `json:"cacheCreationTokens,omitempty"`
	// CostUSD is what the provider says the invocation cost.
	CostUSD float64 `json:"costUsd,omitempty"`
}

// Add accumulates another invocation's usage. A nil addend changes nothing,
// so a run mixing reporting and non-reporting providers still totals what
// was actually reported rather than refusing to total at all.
func (u *Usage) Add(other *Usage) {
	if other == nil {
		return
	}
	u.InputTokens += other.InputTokens
	u.OutputTokens += other.OutputTokens
	u.CacheReadTokens += other.CacheReadTokens
	u.CacheCreationTokens += other.CacheCreationTokens
	u.CostUSD += other.CostUSD
}

// TotalUsage sums every stage's reported usage. Stages that reported
// nothing contribute nothing; the total is of what was measured, and the
// absence of a number is not treated as a zero.
func (s *RunState) TotalUsage() Usage {
	var total Usage
	for _, record := range s.Stages {
		total.Add(record.Usage)
	}
	return total
}
