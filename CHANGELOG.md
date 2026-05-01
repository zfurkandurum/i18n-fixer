# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.10] - 2026-05-01

### Fixed

- **ICU plural objects no longer reported as missing keys.** When a translation
  file stores a plural form as a nested object whose leaf keys are `zero`/`one`/
  `two`/`few`/`many`/`other` (the standard payload shape for easy_localization
  `.plural()`, i18next plural objects, Rails YAML pluralization, etc.), the
  parser now collapses that object into a single entry under the parent key
  using the `other` form as the canonical value. Previously the parser emitted
  three flat sub-keys (e.g. `a.b.zero`, `a.b.one`, `a.b.other`) and never
  registered the parent (`a.b`), so callers like `t('a.b', {count})` or
  `'a.b'.plural(count)` were incorrectly flagged as missing in every locale.
  Non-plural nested objects, mixed-key objects (`{small, other}`), and objects
  missing the mandatory `other` form are unaffected.
- **Interpolated keys no longer drop their static prefix.** When a captured
  i18n key contains runtime interpolation — JS/TS/Dart/Kotlin/PHP `${var}`,
  Dart bare `$var`, Ruby `#{var}`, Swift `\(var)`, Python f-string `{var}`,
  or string concat `"a.b." + var` — the scanner now extracts the static
  prefix up to the last separator before the interpolation and registers it
  as a dynamic prefix. The unused-key analyzer treats every defined key
  starting with that prefix as used, eliminating bulk false-positive
  "unused" reports for runtime-resolved key spaces (typical pattern: enum
  `displayName` getters that build keys from the enum value's name).
- **`IsDynamicKey` recognises three additional interpolation flavours:**
  Ruby `#{`, Swift `\(`, and Python `{var}` (as well as any stray `{` —
  legitimately-typed static i18n keys never contain a brace). Combined with
  the existing `${`/`$`/`+`/backtick/camelCase heuristics, the dynamic-key
  detector now covers the dominant interpolation syntaxes across the
  framework presets bundled with this tool.

### Added

- **`dynamicPrefixPatterns` regex on `flutter-easy-localization` and
  `flutter-getx` presets.** These two presets use a strict captured-key
  character class (`[a-zA-Z_][a-zA-Z0-9_.$]*`) that fails to match any call
  containing `${...}` interpolation, so the scanner-level prefix inference
  alone could not help them. The new patterns explicitly capture the static
  prefix from interpolated `.tr()`/`.plural()`/`.trArgs()` calls; the
  existing `ScanDynamicPrefixes` machinery wires them into the unused-key
  filter. Other framework presets (i18next, vue-i18n, ngx-translate,
  rails-erb, django, react-i18next, etc.) use a loose `[^'"]+` capture
  class and benefit from the scanner-level fix automatically without any
  preset edits.

### Tests

- `TestParseJSONPlural` and `TestParseYAMLPlural` cover the plural-collapse
  cases (full ICU set, one+other only, mixed non-category siblings, missing
  `other`, non-string values).
- `TestIsPluralObject` table-driven sub-cases.
- `TestExtractStaticPrefix` table-driven sub-cases for every supported
  interpolation flavour.
- `TestScanKeyUsageInfersPrefixes` exercises the end-to-end flow: capture
  → dynamic warning → inferred prefix.
- `TestFindMissingKeysPluralRegression` and `TestFindUnusedKeysDynamicPrefixRegression`
  guard against re-introducing either bug at the analyzer level.

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
