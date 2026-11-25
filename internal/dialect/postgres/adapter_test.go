package postgres_test

import (
	"testing"

	"github.com/phantom-mbx1/schemapulse/internal/dialect"
	"github.com/phantom-mbx1/schemapulse/internal/dialect/postgres"
)

func TestParseStatementsIgnoresQuotesAndComments(t *testing.T) {
	sql := `-- comment with semicolon ;
SELECT 1;

INSERT INTO notes (body) VALUES ('semi;colon inside string');

/* block comment
   with ; semicolon */
CREATE INDEX CONCURRENTLY idx_notes_body ON notes (body);
`

	adapter := postgres.New()
	stmts, err := adapter.ParseStatements("0001_notes.sql", sql)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(stmts) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(stmts))
	}
	if stmts[2].Kind != dialect.KindCreateIndex || !stmts[2].Concurrent {
		t.Fatalf("expected statement 3 to be concurrent CREATE INDEX, got %+v", stmts[2])
	}
}

func TestParseStatementsUnterminatedQuoteError(t *testing.T) {
	adapter := postgres.New()
	_, err := adapter.ParseStatements("bad.sql", "INSERT INTO t VALUES ('unterminated);")
	if err == nil {
		t.Fatalf("expected error for unterminated single quote")
	}
}
