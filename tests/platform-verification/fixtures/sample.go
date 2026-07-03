// Intentional test fixture for tests/platform-verification/ — NOT real framework code.
// Deliberately violates go-backend.mdc / go-backend.instructions.md / architecture.mdc rules so a
// tool that actually loaded those rules has something obvious to flag:
//   - interface{} instead of a typed interface
//   - raw SQL string concatenation instead of parameterized queries
//   - no explicit timeout on the HTTP call
//   - error silently swallowed (ignored return value)
//   - domain-ish function directly depends on database/sql (no repository interface)

package sample

import (
	"database/sql"
	"fmt"
	"net/http"
)

func GetUser(db *sql.DB, userID string) interface{} {
	query := "SELECT * FROM users WHERE id = " + userID
	row, _ := db.Query(query)

	resp, _ := http.Get("https://internal-api/enrich?user=" + userID)
	defer resp.Body.Close()

	fmt.Println("fetched user " + userID)
	return row
}
