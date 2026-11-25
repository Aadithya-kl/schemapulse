package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/phantom-mbx1/schemapulse/internal/config"
	"github.com/phantom-mbx1/schemapulse/internal/dialect"
)

type pg005NotNullColumn struct{}

var addNotNullRe = regexp.MustCompile(`(?i)\bADD\s+(?:COLUMN\s+)?(?:IF\s+NOT\s+EXISTS\s+)?([A-Za-z0-9_."]+)\s+([^,;]+?)\s+NOT\s+NULL\b`)

func NewPG005NotNullColumn() Rule { return &pg005NotNullColumn{} }

func (r *pg005NotNullColumn) ID() string { return "PG005" }

func (r *pg005NotNullColumn) Title() string {
	return "Adding NOT NULL column without a constant DEFAULT"
}

func (r *pg005NotNullColumn) DefaultSeverity() config.Severity { return config.SeverityMedium }

func (r *pg005NotNullColumn) Evaluate(stmt dialect.Statement) []Finding {
	if stmt.Kind != dialect.KindAlterTable {
		return nil
	}
	m := addNotNullRe.FindStringSubmatch(stmt.Normalized)
	if len(m) < 3 {
		return nil
	}
	col := strings.Trim(m[1], "\"")
	colDef := strings.ToUpper(m[2])
	hasVolatile := strings.Contains(colDef, "RANDOM(") || strings.Contains(colDef, "CLOCK_TIMESTAMP(") || strings.Contains(colDef, "NOW()")
	if stmt.HasDefault && !hasVolatile {
		return nil
	}
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
			Message:     fmt.Sprintf("Adding NOT NULL column %q without a safe constant DEFAULT may fail on existing rows or trigger a rewrite.", obj),
			Remediation: "Add the column with a constant DEFAULT (or nullable first, backfill, then add CHECK NOT VALID + VALIDATE CONSTRAINT).",
			Statement:   stmt.RawSQL,
		},
	}
}
