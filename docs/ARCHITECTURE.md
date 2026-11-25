# SchemaPulse Architecture

## Pipeline Overview

```text
cmd/schemapulse/main.go (CLI & exit codes)
  ├── config.Load() / WriteStarter()        -> internal/config
  ├── migration.Scan()                      -> internal/migration
  ├── postgres.New().ParseStatements()      -> internal/dialect & internal/dialect/postgres
  ├── rules.DefaultRules()                  -> internal/rules
  ├── analyzer.NewWithOptions().Analyze()    -> internal/analyzer
  └── reporter.New(format).Render()         -> internal/reporter (text, json, markdown)
```

## Key Design Principles

1. **Separation of Concerns:** File discovery, SQL statement splitting/classification, rule evaluation, threshold policy, and output formatting live in separate packages.
2. **Extensible Interfaces:**
   - `dialect.Adapter` decouples SQL statement parsing from rule evaluation so new SQL dialects can be added cleanly.
   - `rules.Rule` isolates each check into its own file and unit test case.
   - `reporter.Reporter` isolates output rendering (`text`, `json`, `markdown`) from analysis logic.
3. **Predictable Exit Codes:**
   - `0`: Completed cleanly with no threshold breach.
   - `1`: Findings or sequence issues breached `--fail-on` (or `inspect` reported sequence issues).
   - `2`: Invalid flags, unreadable path, invalid config, or SQL parse error.