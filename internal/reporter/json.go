package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/phantom-mbx1/schemapulse/internal/analyzer"
)

type JSONReporter struct{}

func (r *JSONReporter) Render(w io.Writer, res *analyzer.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(res); err != nil {
		return fmt.Errorf("failed to encode JSON report: %w", err)
	}
	return nil
}
