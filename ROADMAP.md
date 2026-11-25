# SchemaPulse Roadmap

## Near-Term (v0.2 – v0.3)
- Add additional PostgreSQL rules (`PG008+`: e.g., `DROP INDEX` without `CONCURRENTLY`, lock-timeout guidance).
- Add SARIF output format (`--format sarif`) for code-scanning tools.
- Provide ready-to-use GitHub Actions workflow snippets for pull request migration checks.

## Medium-Term (v0.4 – v0.6)
- Add optional migration history / schema drift checks.
- Add modular SQLite and MySQL dialect adapters implementing `dialect.Adapter`.
- Publish cross-platform release binaries (`linux`, `darwin`, `windows` across `amd64` and `arm64`) on GitHub Releases.

## Long-Term (v1.0+)
- Stabilize rule configuration schema, custom rule conventions, and exit code contracts.
- Support migration runner directory conventions across popular SQL migration tools.