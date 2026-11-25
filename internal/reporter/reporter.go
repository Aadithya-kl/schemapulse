package reporter

import (
	"fmt"
	"io"

	"github.com/phantom-mbx1/schemapulse/internal/analyzer"
)

type Reporter interface {
	Render(w io.Writer, res *analyzer.Result) error
}

func New(format string) (Reporter, error) {
	switch format {
	case "", "text", "cli":
		return &TextReporter{}, nil
	case "json":
		return &JSONReporter{}, nil
	case "markdown", "md":
		return &MarkdownReporter{}, nil
	default:
		return nil, fmt.Errorf("unsupported output format %q (supported: text, json, markdown)", format)
	}
}
