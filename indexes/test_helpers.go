package indexes

import (
	"path/filepath"
	"testing"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a temporary BadgerDB for testing.
// The database is automatically cleaned up when the test completes.
func setupTestDB(t *testing.T) *badger.DB {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "badger")

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	require.NoError(t, err, "Failed to open test database")

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close test database: %v", err)
		}
	})

	return db
}

// setupTestIndexer creates a new Indexer with a temporary database.
func setupTestIndexer(t *testing.T) *Indexer {
	t.Helper()

	db := setupTestDB(t)
	return NewIndexer(db)
}
