# Contributing to SchemaPulse

Thank you for your interest in contributing to **SchemaPulse**! Contributions of any size—new PostgreSQL safety rules, parser improvements, test fixtures, documentation fixes, or bug reports—are welcome.

## Local Development Setup

1. Install **Go** (`1.23+` recommended; tested with Go `1.27`).
2. Fork and clone your fork of `schemapulse`:
   ```bash
   git clone https://github.com/<your-username>/schemapulse.git
   cd schemapulse
   ```
3. Run formatting, static analysis, and unit tests:
   ```bash
   go fmt ./...
   go vet ./...
   go test -v ./...
   ```

## Repository Layout & Independent Modules

- `internal/rules/`: Each rule has its own file (`pg001_*.go` .. `pg007_*.go`) implementing `rules.Rule` (`ID()`, `Title()`, `DefaultSeverity()`, `Evaluate(dialect.Statement) []Finding`) and is registered in `internal/rules/rule.go`.
- `internal/dialect/`: Defines `dialect.Adapter` and `dialect.Statement` in `dialect.go`, with dialect implementations under `internal/dialect/<dialect>/`.
- `internal/migration/`: Handles `.sql` discovery and migration sequence/naming checks in `scanner.go`.
- `internal/config/`: Handles default config, JSON loading, validation, severity overrides, and `schemapulse init`.
- `internal/analyzer/`: Coordinates scanning, parsing, rule filtering (`--only` / `--skip`), severity summaries, and `--fail-on` thresholds.
- `internal/reporter/`: Formats output for `text`, `json`, and `markdown`.
- `cmd/schemapulse/`: CLI entrypoint, subcommands (`check`, `inspect`, `rules`, `init`, `version`, `help`), and exit codes (`0`, `1`, `2`).

## How to Add a New Rule

1. Create `internal/rules/pg008_<name>.go` implementing `rules.Rule`.
2. Add the constructor to `DefaultRules()` in `internal/rules/rule.go`.
3. Add unit tests in `internal/rules/rules_test.go` covering both risky and safe SQL statements.
4. Update the rule table in `README.md` and note the change in `CHANGELOG.md`.
5. Run `go fmt ./...`, `go vet ./...`, and `go test -v ./...`.

## Pull Request Checklist

- Keep PRs focused on a single issue or feature.
- Add or update tests for new rules, flags, or bug fixes.
- Ensure no compiled binaries (`schemapulse`, `schemapulse.exe`, `coverage.out`) are committed.
- Summarize what changed, why, and the test commands you ran.