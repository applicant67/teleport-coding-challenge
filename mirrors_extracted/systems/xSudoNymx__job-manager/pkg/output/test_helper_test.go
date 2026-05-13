package output

import (
	"path/filepath"
	"testing"
)

// newTestLog creates a log in a temp directory.
func newTestLog(t *testing.T) *Log {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	log, err := New(Config{Path: path})
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Cleanup(func() {
		log.Close()
	})

	return log
}
