package migration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/phantom-mbx1/schemapulse/internal/migration"
)

func TestScanDetectsDuplicateSequenceAndValidOrdering(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"0001_first.sql":  "SELECT 1;",
		"0002_second.sql": "SELECT 1;",
		"0002_dup.sql":    "SELECT 1;",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("failed writing test fixture %s: %v", name, err)
		}
	}

	res, err := migration.Scan(dir)
	if err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
	if len(res.Files) != 3 {
		t.Fatalf("expected 3 files scanned, got %d", len(res.Files))
	}
	if len(res.Issues) != 1 {
		t.Fatalf("expected 1 duplicate sequence issue, got %+v", res.Issues)
	}
}
