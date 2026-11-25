package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/phantom-mbx1/schemapulse/internal/config"
	"github.com/phantom-mbx1/schemapulse/internal/dialect"
)

type pg003DropColumn struct{}

var dropColExtractRe = regexp.MustCompile(`(?i)\bDROP\s+(?:COLUMN\s+)?(?:IF\s+EXISTS\s+)?([A-Za-z0-9_."]+)`)

func NewPG003DropColumn() Rule { return &pg003DropColumn{} }

func (r *pg003DropColumn) ID() string { return "PG003" }

func (r *pg003DropColumn) Title() string {
	return "Destructive DROP COLUMN"
}

func (r *pg003DropColumn) DefaultSeverity() config.Severity { return config.SeverityMedium }

func (r *pg003DropColumn) Evaluate(stmt dialect.Statement) []Finding {
	if stmt.Kind != dialect.KindAlterTable {
		return nil
	}
	m := dropColExtractRe.FindStringSubmatch(stmt.Normalized)
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
			Message:     fmt.Sprintf("Dropping column %q permanently discards column data and may break code still reading it during rolling deploys.", obj),
			Remediation: "Stop reading/writing the column in application code first, deploy, and drop the column in a subsequent migration.",
			Statement:   stmt.RawSQL,
		},
	}
}
