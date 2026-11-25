# Security Policy

## Supported Versions

Active security support applies to the latest `0.x` release and `main` branch of **SchemaPulse**.

## Security & Privacy Practices

- `SchemaPulse` performs static analysis on SQL migration files locally and does not require database credentials or network access by default.
- Never include real production credentials, API keys, or sensitive customer records in migration test fixtures or issue reports.

## Reporting a Vulnerability

If you discover a security issue (such as unsafe path handling or parser crashes on malformed input), please **do not** disclose exploit details publicly in an issue before maintainers have had a chance to investigate. Contact the repository maintainers privately first so a fix can be prepared and coordinated responsibly.