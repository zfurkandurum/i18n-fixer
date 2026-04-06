# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-04-06

### Added

- Locale completeness percentage — shows translation coverage per locale
- Duplicate key detection — finds same key with conflicting values across files
- Key naming convention lint — validates UPPER_SNAKE, lower.dot, camelCase, kebab-case with auto-detection
- CLI flags: `--no-duplicates`, `--no-naming`, `--no-completeness`, `--key-convention`
- Overall completeness percentage in summary

### Fixed

- HTML template scanning: Angular `{{ 'KEY' | translate }}` pipe syntax now correctly extracts keys (reduced false unused key reports by ~73%)

## [0.1.0] - 2026-04-06

### Added

- Initial release
- Framework auto-detection (React, Vue, Angular, Svelte, Next.js, Nuxt, Ember, Flutter, iOS, Android, React Native)
- 13 built-in framework presets
- Custom preset support via JSON files
- Missing translation key detection
- Unused translation key detection
- Hardcoded string detection with suggested i18n keys
- Dynamic key detection with manual review warnings
- Console output with formatted summary table
- JSON output for CI/CD integration
- AI prompt output for automated fixing via Claude/GPT
- i18n file parsers: JSON (flat/nested), YAML, Android XML, iOS .strings, Flutter .arb
- Parallel source file scanning via goroutines
- Configuration via `.i18n-fixer.json` with directory walk-up
- npm distribution package with platform-specific binaries
- Exit code 1 when issues found (CI-friendly)
