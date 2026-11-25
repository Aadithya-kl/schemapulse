package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/phantom-mbx1/schemapulse/internal/analyzer"
)

type MarkdownReporter struct{}

func (r *MarkdownReporter) Render(w io.Writer, res *analyzer.Result) error {
	fmt.Fprintf(w, "## SchemaPulse Migration Report (`%s`)\n\n", res.Dialect)
	fmt.Fprintf(w, "- **Target:** `%s`\n", res.Target)
	fmt.Fprintf(w, "- **Files Analyzed:** %d\n", res.FilesAnalyzed)
	fmt.Fprintf(w, "- **Statements Scanned:** %d\n", res.StatementsTotal)
	fmt.Fprintf(w, "- **Findings:** %d\n", len(res.Findings))
	fmt.Fprintf(w, "- **Sequence Issues:** %d\n\n", len(res.SequenceIssues))

	if len(res.SequenceIssues) > 0 {
		fmt.Fprintln(w, "### Migration Sequence Issues")
		for _, issue := range res.SequenceIssues {
			fmt.Fprintf(w, "- `%s`: %s\n", issue.File, issue.Message)
		}
		fmt.Fprintln(w)
	}

	if len(res.Findings) > 0 {
		fmt.Fprintln(w, "### Findings")
		fmt.Fprintln(w, "| Rule | Severity | Location | Issue | Suggested Fix |")
		fmt.Fprintln(w, "| --- | --- | --- | --- | --- |")
		for _, f := range res.Findings {
			loc := fmt.Sprintf("`%s:%d`", f.File, f.Line)
			msg := strings.ReplaceAll(f.Message, "|", "\\|")
			fix := strings.ReplaceAll(f.Remediation, "|", "\\|")
			fmt.Fprintf(w, "| `%s` | **%s** | %s | %s | %s |\n",
				f.RuleID, strings.ToUpper(string(f.Severity)), loc, msg, fix)
		}
		fmt.Fprintln(w)
	} else {
		fmt.Fprintln(w, "No rule findings detected.")
	}

	return nil
}
