package rules

import (
	"fmt"

	"github.com/phantom-mbx1/schemapulse/internal/config"
	"github.com/phantom-mbx1/schemapulse/internal/dialect"
)

type pg002DropTable struct{}

func NewPG002DropTable() Rule { return &pg002DropTable{} }

func (r *pg002DropTable) ID() string { return "PG002" }

func (r *pg002DropTable) Title() string {
	return "Destructive DROP TABLE or DROP SCHEMA"
}

func (r *pg002DropTable) DefaultSeverity() config.Severity { return config.SeverityHigh }

func (r *pg002DropTable) Evaluate(stmt dialect.Statement) []Finding {
	if stmt.Kind != dialect.KindDropTable && stmt.Kind != dialect.KindDropSchema {
		return nil
	}
	kindLabel := "TABLE"
	if stmt.Kind == dialect.KindDropSchema {
		kindLabel = "SCHEMA"
	}
	target := stmt.ObjectName
	if target == "" {
		target = "target object"
	}
	return []Finding{
		{
			RuleID:      r.ID(),
			Title:       r.Title(),
			Severity:    r.DefaultSeverity(),
			File:        stmt.File,
			Line:        stmt.Line,
			Object:      target,
			Message:     fmt.Sprintf("DROP %s on %q permanently removes data and can break running application queries.", kindLabel, target),
			Remediation: "Verify backups, ensure no active code references the object, and consider multi-step deprecation before dropping.",
			Statement:   stmt.RawSQL,
		},
	}
}
