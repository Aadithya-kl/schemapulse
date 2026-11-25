package rules_test

import (
	"testing"

	"github.com/phantom-mbx1/schemapulse/internal/dialect"
	"github.com/phantom-mbx1/schemapulse/internal/dialect/postgres"
	"github.com/phantom-mbx1/schemapulse/internal/rules"
)

func TestPG007SetNotNull(t *testing.T) {
	adapter := postgres.New()
	stmts, err := adapter.ParseStatements("0010_nn.sql", "ALTER TABLE orders ALTER COLUMN customer_id SET NOT NULL;")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(stmts) != 1 || stmts[0].Kind != dialect.KindAlterTable {
		t.Fatalf("unexpected parsed statements: %+v", stmts)
	}

	r := rules.NewPG007SetNotNull()
	findings := r.Evaluate(stmts[0])
	if len(findings) != 1 || findings[0].RuleID != "PG007" {
		t.Fatalf("expected 1 PG007 finding, got %+v", findings)
	}
}
