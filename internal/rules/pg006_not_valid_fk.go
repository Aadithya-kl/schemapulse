package rules

import (
	"fmt"
	"strings"

	"github.com/phantom-mbx1/schemapulse/internal/config"
	"github.com/phantom-mbx1/schemapulse/internal/dialect"
)

type pg006NotValidForeignKey struct{}

func NewPG006NotValidForeignKey() Rule { return &pg006NotValidForeignKey{} }

func (r *pg006NotValidForeignKey) ID() string { return "PG006" }

func (r *pg006NotValidForeignKey) Title() string {
	return "Adding FOREIGN KEY without NOT VALID"
}

func (r *pg006NotValidForeignKey) DefaultSeverity() config.Severity { return config.SeverityHigh }

func (r *pg006NotValidForeignKey) Evaluate(stmt dialect.Statement) []Finding {
	if stmt.Kind != dialect.KindAlterTable {
		return nil
	}
	upper := strings.ToUpper(stmt.Normalized)
	if !strings.Contains(upper, "FOREIGN KEY") {
		return nil
	}
	if stmt.HasNotValid {
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
			Message:     fmt.Sprintf("Adding a FOREIGN KEY on %q without NOT VALID scans all existing rows while holding a lock.", target),
			Remediation: "Add the foreign key constraint with NOT VALID first, then run ALTER TABLE ... VALIDATE CONSTRAINT ... in a separate step.",
			Statement:   stmt.RawSQL,
		},
	}
}
