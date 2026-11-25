package analyzer_test

import (
	"path/filepath"
	"testing"

	"github.com/phantom-mbx1/schemapulse/internal/analyzer"
	"github.com/phantom-mbx1/schemapulse/internal/config"
	"github.com/phantom-mbx1/schemapulse/internal/dialect/postgres"
	"github.com/phantom-mbx1/schemapulse/internal/rules"
)

func TestAnalyzeRiskyAndSafeMigrations(t *testing.T) {
	eng := analyzer.New(postgres.New(), rules.DefaultRules(), config.Default())

	riskyPath := filepath.Join("..", "..", "testdata", "migrations", "risky", "0001_risky_changes.sql")
	riskyRes, err := eng.Analyze(riskyPath)
	if err != nil {
		t.Fatalf("unexpected error analyzing risky migration: %v", err)
	}
	if len(riskyRes.Findings) < 6 {
		t.Fatalf("expected at least 6 findings for risky migration, got %d", len(riskyRes.Findings))
	}
	if !riskyRes.ThresholdBreach {
		t.Fatalf("expected ThresholdBreach=true for risky migration with failOn=high")
	}

	safePath := filepath.Join("..", "..", "testdata", "migrations", "safe", "0001_safe_changes.sql")
	safeRes, err := eng.Analyze(safePath)
	if err != nil {
		t.Fatalf("unexpected error analyzing safe migration: %v", err)
	}
	if len(safeRes.Findings) != 0 {
		t.Fatalf("expected 0 findings for safe migration, got %d: %+v", len(safeRes.Findings), safeRes.Findings)
	}
	if safeRes.ThresholdBreach {
		t.Fatalf("expected ThresholdBreach=false for safe migration")
	}
}
