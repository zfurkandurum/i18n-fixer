# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.9] - 2026-04-26

### Fixed

- **Strings inside `return`, `=>`, and `throw` are now detected.** Previously the scanner only matched explicit widget/HTML-attribute contexts (`Text(...)`, `placeholder="..."`, etc.), so hardcoded user-facing literals hidden inside enum extensions, switch arrow arms, switch/case bodies, error throws, and arrow-getter returns slipped through silently. All 31 framework presets now include language-appropriate `return`/`=>`/`throw`/`raise` patterns. Extracted strings still pass through the existing exclusion chain — translation keys (`'common.cancel'.tr()`), `$variable` interpolation, asset paths, etc. remain ignored, so this change adds coverage without flooding reports with false positives.
- **Dotted-identifier exclusion now anchored.** All four Flutter presets had `^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z0-9_]+)+` (no end anchor), which silently swallowed sentences that happened to start with a dotted token (`"e.g. Jane Doe"`, `"i.e. ..."`). The `$` end-anchor restores the intended "entire string is a translation key" semantics.

### Added

- **Expanded named-parameter coverage per language.**
  - Dart (4 presets): `label:`, `subtitle:`, `helperText:`, `errorText:`, `placeholder:`, `actionText:`, `description:`, `text:`, `throw '...'`, `Exception('...')`.
  - Swift (4 iOS presets): `subtitle:`, `description:`, `accessibilityLabel:`.
  - Kotlin (Android, Compose): `description=`, `supportingText=`, `throw *Exception(...)`.
  - JS/TS (React, Next, Lingui, etc.): JSX attributes `subtitle`, `description`, `helperText`, `errorText`, `caption`, `heading`, `header`, `footer`; `throw new Error('...')` where missing.
  - Python (Django), Ruby (Rails), PHP (Laravel), Go (go-i18n): `return` literals; Ruby also gets `raise '...'`.

### Tests

- New `internal/scanner/testdata/sample_dart_enums.dart` fixture and `TestScanHardcodedDartEnumsAndReturns` covering arrow switch arms, classic `case … : return …;`, `throw '…'`, `throw Exception('…')`, and the full extended named-parameter list. Translation keys (`'common.cancel'.tr()`, `'foo.bar.$name'.tr()`) and interpolated values remain excluded.

## [0.3.2] - 2026-04-07

### Fixed

- Windows compatibility: preset loading now uses forward slashes for `embed.FS` paths
- CI matrix updated to Go 1.26 to match `go.mod` minimum requirement

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
