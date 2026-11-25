package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/phantom-mbx1/schemapulse/internal/config"
	"github.com/phantom-mbx1/schemapulse/internal/dialect"
)

type pg004AlterColumnType struct{}

var alterTypeColRe = regexp.MustCompile(`(?i)\bALTER\s+(?:COLUMN\s+)?([A-Za-z0-9_."]+)\s+(?:SET\s+DATA\s+)?TYPE\b`)

func NewPG004AlterColumnType() Rule { return &pg004AlterColumnType{} }

func (r *pg004AlterColumnType) ID() string { return "PG004" }

func (r *pg004AlterColumnType) Title() string {
	return "ALTER COLUMN TYPE may rewrite table or lock rows"
}

func (r *pg004AlterColumnType) DefaultSeverity() config.Severity { return config.SeverityHigh }

func (r *pg004AlterColumnType) Evaluate(stmt dialect.Statement) []Finding {
	if stmt.Kind != dialect.KindAlterTable {
		return nil
	}
	m := alterTypeColRe.FindStringSubmatch(stmt.Normalized)
	if len(m) < 2 {
		return nil
	}
	col := strings.Trim(m[1], "\"")
	obj := col
	if stmt.ObjectName != "" {
		obj = stmt.ObjectName + "." + col
	}
	return []Finding{
		{
			RuleID:      r.ID(),
			Title:       r.Title(),
			Severity:    r.DefaultSeverity(),
			File:        stmt.File,
			Line:        stmt.Line,
			Object:      obj,
			Message:     fmt.Sprintf("Altering column type on %q can require a full table rewrite under an ACCESS EXCLUSIVE lock.", obj),
			Remediation: "Prefer expand-and-contract (add new column, backfill in batches, dual-write, swap, drop old) for large tables.",
			Statement:   stmt.RawSQL,
		},
	}
}
