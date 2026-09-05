package memory

// Run identity.
//
// The executor writes no run ID — it has never needed one, since a run is
// identified by the workspace it lives in. This epic must not modify the
// executor, so the identity is derived from what run state already carries.
//
// Feature plus start time is the right pair. It is stable across a resume
// (StartedAt is set once, when the run is created, and survives every
// checkpoint), which is exactly the property re-ingest needs: ingesting a
// run, resuming it, and ingesting again must update one row rather than
// making two. And it does not collide across features, or across two runs
// of the same feature started at different times.

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// RunID derives a stable identifier for one run.
func RunID(feature string, startedAt time.Time) string {
	sum := sha256.Sum256([]byte(feature + "@" + startedAt.UTC().Format(time.RFC3339Nano)))
	return hex.EncodeToString(sum[:12])
}
