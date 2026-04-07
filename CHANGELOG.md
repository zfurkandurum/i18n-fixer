# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.1] - 2026-04-07

### Fixed

- **Template expression false positives** — added `{{`, `{`, `{%`, `<%=` exclusions to all 31 presets so Angular/Vue/Svelte/Django/Rails template syntax is never flagged as a hardcoded string
- **Test file scanning** — all presets now exclude test/spec files (`*.spec.ts`, `*_test.dart`, `*_spec.rb`, `*Test.kt`, etc.) to prevent fixture strings from polluting reports
- **Flutter dot-notation false positives** — `"key.name".tr()` inside `Text()` calls no longer reported as hardcoded string
- **flutter-easy-localization: `tr('key')` function-call pattern** — added `tr('key')` and `translate('key')` patterns alongside the existing `'key'.tr()` extension style

## [0.3.0] - 2026-04-07

### Added

- **18 new built-in presets** (total: 31)
  - JavaScript/TypeScript: `i18next`, `next-i18next`, `next-translate`, `lingui`, `typesafe-i18n`, `inlang-paraglide`
  - Flutter: `flutter-easy-localization`, `flutter-getx`, `flutter-slang`
  - iOS/Swift: `ios-swiftui`, `ios-swiftgen`, `ios-xcstrings`
  - Android: `android-compose`
  - React Native: `react-native-i18n-js`
  - Backend/Templates: `django`, `rails-erb`, `laravel-blade`, `go-i18n`
- **`i18n-ignore` annotation** — add `// i18n-ignore` (or any comment style) to a line to suppress false positives
- **`.xcstrings` parser** — Apple String Catalog format (Xcode 15+) with full plural variant support; auto-detected by file extension
- **GitHub Actions integration** — `action.yml` for `uses: zfurkandurum/i18n-fixer@v1`
- **Pre-commit hook** — `.pre-commit-hooks.yaml` for pre-commit framework integration
- **HTML/template file scanning** — `.html`, `.tmpl`, `.erb`, `.blade.php` extensions covered by new presets
- **`xcstrings`** added as a valid `i18nFileFormat` value in custom presets

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
