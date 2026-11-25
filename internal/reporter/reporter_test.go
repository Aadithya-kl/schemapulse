package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/phantom-mbx1/schemapulse/internal/analyzer"
	"github.com/phantom-mbx1/schemapulse/internal/reporter"
)

func TestReportersRenderValidOutput(t *testing.T) {
	res := &analyzer.Result{
		Target:          "migrations",
		Dialect:         "postgres",
		FilesAnalyzed:   1,
		StatementsTotal: 2,
		FailOn:          "high",
	}

	var jsonBuf bytes.Buffer
	jr, err := reporter.New("json")
	if err != nil {
		t.Fatalf("unexpected error creating json reporter: %v", err)
	}
	if err := jr.Render(&jsonBuf, res); err != nil {
		t.Fatalf("unexpected error rendering json: %v", err)
	}
	var decoded analyzer.Result
	if err := json.Unmarshal(jsonBuf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	var mdBuf bytes.Buffer
	mr, err := reporter.New("markdown")
	if err != nil {
		t.Fatalf("unexpected error creating markdown reporter: %v", err)
	}
	if err := mr.Render(&mdBuf, res); err != nil {
		t.Fatalf("unexpected error rendering markdown: %v", err)
	}
	if !strings.Contains(mdBuf.String(), "## SchemaPulse Migration Report") {
		t.Fatalf("expected markdown header, got %q", mdBuf.String())
	}
}
