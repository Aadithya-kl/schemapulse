package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/phantom-mbx1/schemapulse/internal/config"
	"github.com/phantom-mbx1/schemapulse/internal/dialect"
)

type pg007SetNotNull struct{}

var setNotNullColRe = regexp.MustCompile(`(?i)\bALTER\s+(?:COLUMN\s+)?([A-Za-z0-9_."]+)\s+SET\s+NOT\s+NULL\b`)

func NewPG007SetNotNull() Rule { return &pg007SetNotNull{} }

func (r *pg007SetNotNull) ID() string { return "PG007" }

func (r *pg007SetNotNull) Title() string {
	return "ALTER COLUMN SET NOT NULL scans table under lock"
}

func (r *pg007SetNotNull) DefaultSeverity() config.Severity { return config.SeverityMedium }

func (r *pg007SetNotNull) Evaluate(stmt dialect.Statement) []Finding {
	if stmt.Kind != dialect.KindAlterTable {
		return nil
	}
	m := setNotNullColRe.FindStringSubmatch(stmt.Normalized)
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
			Message:     fmt.Sprintf("Setting NOT NULL on existing column %q requires scanning all rows while holding an ACCESS EXCLUSIVE lock.", obj),
			Remediation: "Add a CHECK constraint with NOT VALID, VALIDATE CONSTRAINT in a separate transaction, or backfill NULLs first to minimize lock duration.",
			Statement:   stmt.RawSQL,
		},
	}
}
