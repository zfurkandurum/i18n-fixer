# Contributing to i18n-fixer

Thank you for your interest in contributing!

## Development Setup

1. Install [Go 1.22+](https://go.dev/dl/)
2. Fork and clone the repository
3. Build: `make build`
4. Test: `make test`
5. Lint: `make lint` (requires [golangci-lint](https://golangci-lint.run/))

## Code Style

- Format with `gofmt -s` (enforced in CI)
- Follow golangci-lint rules
- Write table-driven tests

## Pull Request Process

1. Create a feature branch (`feat/xxx` or `fix/xxx`)
2. Write or update tests
3. Ensure `make all` passes (fmt + vet + test + build)
4. Update CHANGELOG.md (add to Unreleased section)
5. Update README.md if adding user-facing features
6. Open a PR with a clear description of what and why

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: add vue-i18n preset`
- `fix: nested key flattening for 4+ depth`
- `docs: update comparison table`
- `test: add xml parser edge cases`
- `chore: update dependencies`

## Adding a New Preset

1. Create `internal/preset/builtin/<name>.json` following the preset schema
2. The preset is auto-loaded via `embed.FS` — no code changes needed for loading
3. Add project markers in the preset JSON for auto-detection
4. Add test fixtures in `internal/detect/testdata/`
5. Update the README.md supported frameworks table

## Project Structure

```
cmd/i18n-fixer/     Entry point
internal/
  cli/              Cobra commands
  config/           Config file loading
  detect/           Framework auto-detection
  preset/           Built-in and custom presets
  parser/           i18n file parsers (JSON, YAML, XML, .strings, .arb)
  scanner/          Source code scanners
  analyzer/         Missing/unused/hardcoded analysis
  reporter/         Output formatters (console, JSON, AI prompt)
  types/            Shared types
```

## Questions?

Open a [Discussion](https://github.com/i18n-fixer/i18n-fixer/discussions) or file an issue.
