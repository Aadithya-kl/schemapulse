package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/phantom-mbx1/schemapulse/internal/analyzer"
)

type TextReporter struct{}

func (r *TextReporter) Render(w io.Writer, res *analyzer.Result) error {
	fmt.Fprintf(w, "SchemaPulse (%s) — Target: %s\n", res.Dialect, res.Target)
	fmt.Fprintf(w, "Files: %d | Statements: %d | Findings: %d (high=%d, medium=%d, low=%d) | Sequence Issues: %d\n",
		res.FilesAnalyzed, res.StatementsTotal, len(res.Findings),
		res.Summary.High, res.Summary.Medium, res.Summary.Low, len(res.SequenceIssues))
	fmt.Fprintln(w, strings.Repeat("-", 72))

	for _, seqErr := range res.SequenceIssues {
		fmt.Fprintf(w, "[SEQUENCE] %s: %s\n", seqErr.File, seqErr.Message)
	}

	for _, f := range res.Findings {
		sevLabel := strings.ToUpper(string(f.Severity))
		loc := fmt.Sprintf("%s:%d", f.File, f.Line)
		if f.Object != "" {
			fmt.Fprintf(w, "[%s] %s (%s) — %s [%s]\n", sevLabel, f.RuleID, loc, f.Title, f.Object)
		} else {
			fmt.Fprintf(w, "[%s] %s (%s) — %s\n", sevLabel, f.RuleID, loc, f.Title)
		}
		fmt.Fprintf(w, "       Why: %s\n", f.Message)
		fmt.Fprintf(w, "       Fix: %s\n", f.Remediation)
	}

	if len(res.Findings) == 0 && len(res.SequenceIssues) == 0 {
		fmt.Fprintln(w, "No migration risks or sequence issues detected.")
	}

	if res.ThresholdBreach {
		fmt.Fprintf(w, "\nResult: FAIL (findings met or exceeded failOn=%q threshold)\n", res.FailOn)
	} else {
		fmt.Fprintf(w, "\nResult: PASS (failOn=%q)\n", res.FailOn)
	}

	return nil
}
