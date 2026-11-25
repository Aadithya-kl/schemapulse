package rules

import (
	"fmt"

	"github.com/phantom-mbx1/schemapulse/internal/config"
	"github.com/phantom-mbx1/schemapulse/internal/dialect"
)

type pg001ConcurrentIndex struct{}

func NewPG001ConcurrentIndex() Rule { return &pg001ConcurrentIndex{} }

func (r *pg001ConcurrentIndex) ID() string { return "PG001" }

func (r *pg001ConcurrentIndex) Title() string {
	return "CREATE INDEX without CONCURRENTLY"
}

func (r *pg001ConcurrentIndex) DefaultSeverity() config.Severity { return config.SeverityHigh }

func (r *pg001ConcurrentIndex) Evaluate(stmt dialect.Statement) []Finding {
	if stmt.Kind != dialect.KindCreateIndex || stmt.Concurrent {
		return nil
	}
	target := stmt.ObjectName
	if target == "" {
		target = "target table"
	}
	return []Finding{
		{
			RuleID:      r.ID(),
			Title:       r.Title(),
			Severity:    r.DefaultSeverity(),
			File:        stmt.File,
			Line:        stmt.Line,
			Object:      target,
			Message:     fmt.Sprintf("Creating an index on %q without CONCURRENTLY acquires a SHARE lock that blocks writes.", target),
			Remediation: "Use CREATE INDEX CONCURRENTLY ... (and run outside a transaction block in your migration runner).",
			Statement:   stmt.RawSQL,
		},
	}
}
