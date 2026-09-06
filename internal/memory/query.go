package memory

// Queries over the episodic store.
//
// Every query here is a named, parameterized statement. There is no
// free-form SQL surface, which costs something real: a question nobody
// anticipated needs a code change. That is accepted deliberately — it keeps
// the schema private so later phases can change it, and it means every
// question the store can answer is documented by `--help` rather than
// discovered by writing SQL against a shape nobody promised to keep.

import (
	"database/sql"
	"fmt"
)

// RunSummary is one run, as anyone asking about history wants to see it.
type RunSummary struct {
	RunID        string `json:"runId"`
	Feature      string `json:"feature"`
	StartedAt    string `json:"startedAt"`
	Complete     bool   `json:"complete"`
	WaitingGate  string `json:"waitingGate,omitempty"`
	InputTokens  int64  `json:"inputTokens"`
	OutputTokens int64  `json:"outputTokens"`
	// CacheReadTokens and CacheCreationTokens are kept because they dominate:
	// a verified real invocation reported 2 input tokens and 58,299 cache
	// creation tokens. A "tokens" figure that omitted them would understate
	// that run by four orders of magnitude.
	CacheReadTokens     int64   `json:"cacheReadTokens"`
	CacheCreationTokens int64   `json:"cacheCreationTokens"`
	CostUSD             float64 `json:"costUsd"`
	Corrections         int     `json:"corrections"`
}

// TotalTokens is every billable token, cache traffic included.
func (r RunSummary) TotalTokens() int64 {
	return r.InputTokens + r.OutputTokens + r.CacheReadTokens + r.CacheCreationTokens
}

// RetryRow is one stage that took more than one attempt.
type RetryRow struct {
	RunID      string `json:"runId"`
	Feature    string `json:"feature"`
	StartedAt  string `json:"startedAt"`
	Stage      string `json:"stage"`
	Agent      string `json:"agent"`
	Iterations int    `json:"iterations"`
}

// AgentCorrections counts how often a human had to fix one agent's output.
type AgentCorrections struct {
	Agent       string `json:"agent"`
	Corrections int    `json:"corrections"`
	LinesAdded  int    `json:"linesAdded"`
	LinesCut    int    `json:"linesRemoved"`
	Runs        int    `json:"runs"`
}

// Runs returns every ingested run, newest first.
func (s *Store) Runs(limit int) ([]RunSummary, error) {
	rows, err := s.db.Query(`
		SELECT r.run_id, r.feature, r.started_at, r.completed, COALESCE(r.waiting_gate, ''),
		       r.input_tokens, r.output_tokens, r.cache_read_tokens, r.cache_creation_tokens,
		       r.cost_usd,
		       (SELECT COUNT(*) FROM corrections c WHERE c.run_id = r.run_id)
		FROM runs r ORDER BY r.started_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("query runs: %w", err)
	}
	return scanRuns(rows)
}

func scanRuns(rows *sql.Rows) ([]RunSummary, error) {
	defer func() { _ = rows.Close() }()
	summaries := make([]RunSummary, 0)
	for rows.Next() {
		var summary RunSummary
		if err := rows.Scan(&summary.RunID, &summary.Feature, &summary.StartedAt, &summary.Complete,
			&summary.WaitingGate, &summary.InputTokens, &summary.OutputTokens,
			&summary.CacheReadTokens, &summary.CacheCreationTokens,
			&summary.CostUSD, &summary.Corrections); err != nil {
			return nil, fmt.Errorf("scan run: %w", err)
		}
		summaries = append(summaries, summary)
	}
	return summaries, rows.Err()
}

// Retries answers L3.5's done-when: which runs needed an agent to go round
// more than a given number of times.
//
// "Retried more than twice" is read as iteration count > 2 — at least three
// attempts. The other reading (three retries after a first attempt) is
// equally defensible, which is why callers print the threshold rather than
// leaving a reader to infer it. An empty agent matches every agent.
func (s *Store) Retries(agent string, moreThan int) ([]RetryRow, error) {
	rows, err := s.db.Query(`
		SELECT s.run_id, r.feature, r.started_at, s.stage_id, COALESCE(s.agent, ''), s.iterations
		FROM stages s JOIN runs r ON r.run_id = s.run_id
		WHERE s.iterations > ? AND (? = '' OR s.agent = ?)
		ORDER BY s.iterations DESC, r.started_at DESC`, moreThan, agent, agent)
	if err != nil {
		return nil, fmt.Errorf("query retries: %w", err)
	}
	return scanRetries(rows)
}

func scanRetries(rows *sql.Rows) ([]RetryRow, error) {
	defer func() { _ = rows.Close() }()
	results := make([]RetryRow, 0)
	for rows.Next() {
		var row RetryRow
		if err := rows.Scan(&row.RunID, &row.Feature, &row.StartedAt, &row.Stage,
			&row.Agent, &row.Iterations); err != nil {
			return nil, fmt.Errorf("scan retry row: %w", err)
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// Corrections ranks agents by how often a human corrected their output —
// the signal epic 85 started collecting and L4.4 will consume.
func (s *Store) Corrections() ([]AgentCorrections, error) {
	rows, err := s.db.Query(`
		SELECT COALESCE(agent, '(unattributed)'), COUNT(*), SUM(added), SUM(removed),
		       COUNT(DISTINCT run_id)
		FROM corrections GROUP BY agent ORDER BY COUNT(*) DESC`)
	if err != nil {
		return nil, fmt.Errorf("query corrections: %w", err)
	}
	return scanCorrections(rows)
}

func scanCorrections(rows *sql.Rows) ([]AgentCorrections, error) {
	defer func() { _ = rows.Close() }()
	results := make([]AgentCorrections, 0)
	for rows.Next() {
		var row AgentCorrections
		if err := rows.Scan(&row.Agent, &row.Corrections, &row.LinesAdded,
			&row.LinesCut, &row.Runs); err != nil {
			return nil, fmt.Errorf("scan correction row: %w", err)
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// CostReported reports whether any run recorded a cost at all. A store full
// of zeroes means the provider reported nothing, not that the runs were
// free — and a cost report that cannot tell those apart is the confident
// wrong number this framework keeps removing.
//
// A nil Store reports false, so a caller rendering a table without a
// database in hand — a test, or a path where the store failed to open —
// gets the cautious answer rather than a panic. Same posture as a nil
// telemetry Session.
func (s *Store) CostReported() (bool, error) {
	if s == nil {
		return false, nil
	}
	var runs int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM runs
		WHERE cost_usd > 0 OR input_tokens > 0 OR output_tokens > 0
		   OR cache_read_tokens > 0 OR cache_creation_tokens > 0`).
		Scan(&runs); err != nil {
		return false, fmt.Errorf("check reported cost: %w", err)
	}
	return runs > 0, nil
}
