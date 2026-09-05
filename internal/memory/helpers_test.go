package memory_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/orieken/loom/internal/memory"
	"github.com/orieken/loom/internal/orchestrator"

	_ "modernc.org/sqlite"
)

// writeWorkspace lays down a run's state and timeline the way the executor
// does, so the readers are tested against real files rather than mocks.
func writeWorkspace(t *testing.T, dir string, state *orchestrator.RunState, events []orchestrator.Event) {
	t.Helper()
	statePath := filepath.Join(dir, orchestrator.RunStateFileName)
	if err := orchestrator.NewStateStore(statePath).Save(state); err != nil {
		t.Fatalf("save state: %v", err)
	}
	timeline := orchestrator.NewTimeline(statePath)
	for _, event := range events {
		if err := timeline.Append(event); err != nil {
			t.Fatalf("append event: %v", err)
		}
	}
}

// openReadOnly reads the store directly, so assertions check what was
// actually written rather than what the writer believes it wrote.
func openReadOnly(t *testing.T, store *memory.Store) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", store.Path())
	if err != nil {
		t.Fatalf("open store for reading: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func assertCounts(t *testing.T, store *memory.Store, want map[string]int) {
	t.Helper()
	db := openReadOnly(t, store)
	for table, expected := range want {
		var got int
		// Table names are literals from the test's own map, never input.
		if err := db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&got); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		if got != expected {
			t.Errorf("%s has %d rows, want %d", table, got, expected)
		}
	}
}

func nullableInt(t *testing.T, store *memory.Store, query string) *int64 {
	t.Helper()
	var value sql.NullInt64
	if err := openReadOnly(t, store).QueryRow(query).Scan(&value); err != nil {
		t.Fatalf("query %q: %v", query, err)
	}
	if !value.Valid {
		return nil
	}
	return &value.Int64
}
