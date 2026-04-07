# i18n-fixer

[![CI](https://github.com/zfurkandurum/i18n-fixer/actions/workflows/ci.yml/badge.svg)](https://github.com/zfurkandurum/i18n-fixer/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/zfurkandurum/i18n-fixer.svg)](https://pkg.go.dev/github.com/zfurkandurum/i18n-fixer)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> Framework-agnostic CLI tool that finds hardcoded strings, missing i18n keys, and unused translations. Generates AI-ready fix prompts. Works with every major frontend, mobile, and backend framework.

## Supported Frameworks

| Category | Frameworks / Libraries |
|----------|----------------------|
| **React** | react-i18next, react-intl (FormatJS) |
| **Next.js** | next-intl, next-i18next, next-translate |
| **Vue** | vue-i18n (Vue 2 & 3) |
| **Nuxt** | @nuxtjs/i18n |
| **Angular** | @angular/localize, @ngx-translate |
| **Svelte** | svelte-i18n |
| **Ember** | ember-intl |
| **JavaScript / TypeScript** | i18next (standalone), LinguiJS, typesafe-i18n, Inlang/Paraglide |
| **Flutter** | intl/ARB, easy_localization, GetX, slang |
| **iOS / Swift** | NSLocalizedString, SwiftUI, SwiftGen, String Catalog (.xcstrings) |
| **Android** | getString(R.string.x) XML resources, Jetpack Compose |
| **React Native** | i18next, react-intl, i18n-js |
| **Django** | gettext, `{% trans %}`, `{% blocktrans %}` |
| **Ruby on Rails** | Rails I18n, `t()`, ERB templates |
| **Laravel** | PHP `__()`, `trans()`, Blade `@lang()` |
| **Go** | go-i18n, html/template |
| **Custom** | Any framework via custom preset JSON |

## Features

- **Auto-detects** your framework — zero config needed
- **31 built-in presets** covering all major frontend, mobile, and backend frameworks
- Finds **missing** translation keys (used in code, absent from i18n files)
- Finds **unused** translation keys (in i18n files, never referenced in code)
- Detects **hardcoded** user-facing strings not wrapped in i18n functions
- **`i18n-ignore` annotation** — suppress false positives with an inline comment
- **Locale completeness %** — see translation coverage per locale at a glance
- **Duplicate key detection** — finds conflicting values for the same key
- **Key naming lint** — enforces UPPER_SNAKE, lower.dot, camelCase, or kebab-case
- Generates **AI-ready fix prompts** — paste into Claude/GPT to auto-fix
- **HTML/template scanning** — Django, ERB, Blade, Go templates
- **`.xcstrings` support** — Apple String Catalog (Xcode 15+) with plural variants
- **Fast** — parallel scanning via goroutines
- **Single binary** — no Node.js, Python, or any runtime needed
- **Extensible** — add any framework via custom preset JSON
- **CI/CD ready** — exit codes, JSON output, GitHub Actions, pre-commit hook

## Quick Start

```bash
# Just run in your project directory — framework auto-detected
i18n-fixer

# Generate AI fix prompt (auto-saves to i18n-fix-prompt.md)
i18n-fixer -f prompt

# JSON report for CI
i18n-fixer -f json -o report.json
```

## Installation

### Homebrew (macOS/Linux)

```bash
brew install i18n-fixer/tap/i18n-fixer
```

### npm / npx (any platform with Node.js)

```bash
# Run without installing
npx i18n-fixer

# Or install globally
npm install -g i18n-fixer
```

### Go

```bash
go install github.com/zfurkandurum/i18n-fixer/cmd/i18n-fixer@latest
```

### Binary Download

Download from [GitHub Releases](https://github.com/zfurkandurum/i18n-fixer/releases).

## Usage

```
i18n-fixer [flags] [path]

Commands:
  run          Scan project for i18n issues (default)
  init         Generate a starter .i18n-fixer.json config
  presets      List available built-in presets
  version      Print version information

Flags:
  -p, --preset <name>          Framework preset name or path to custom JSON
  -f, --format <type>          Output format: console, json, prompt (default: console)
  -o, --output <path>          Write report to file
      --no-hardcoded           Skip hardcoded string detection
      --no-missing             Skip missing key detection
      --no-unused              Skip unused key detection
      --no-duplicates          Skip duplicate key detection
      --no-naming              Skip key naming convention lint
      --no-completeness        Skip locale completeness analysis
      --key-convention <style> Enforce naming: UPPER_SNAKE, lower.dot, camelCase, kebab-case
      --default-locale <code>  Only check missing keys against this locale
      --strict-unused          Disable dynamic key heuristic exclusion
      --ignore <pattern>       Additional glob patterns to ignore (repeatable)
      --verbose                Show detailed scanning progress
      --no-color               Disable colored output
```

### Examples

```bash
# Auto-detect and scan current directory
i18n-fixer

# Scan a specific directory
i18n-fixer ./frontend

# Use a specific preset
i18n-fixer --preset react-i18next
i18n-fixer --preset flutter-easy-localization
i18n-fixer --preset django

# Generate AI prompt (auto-saves to i18n-fix-prompt.md)
i18n-fixer -f prompt

# JSON output for CI pipelines
i18n-fixer -f json -o report.json

# Only check for missing keys
i18n-fixer --no-hardcoded --no-unused

# Enforce key naming convention
i18n-fixer --key-convention lower.dot

# Check against a single locale
i18n-fixer --default-locale en

# Use custom preset for unsupported framework
i18n-fixer --preset ./my-preset.json

# List all available presets
i18n-fixer presets
```

## Translation File Location

i18n-fixer automatically finds your translation files — no configuration needed. It searches the entire project for files in any of these directory names, at any depth:

```
locales/    locale/    i18n/    lang/    translations/    messages/
```

For example, all of these are discovered automatically:

```
src/locales/en.json
src/assets/i18n/tr.json
apps/web/src/lang/fr.json          ← monorepo
packages/ui/locales/de.json        ← monorepo
public/locales/en/common.json      ← next-i18next namespace style
```

**Platform-specific conventions** (always auto-detected):

| Platform | Where i18n-fixer looks |
|----------|----------------------|
| Android | `**/res/values*/strings.xml` |
| iOS | `**/*.lproj/Localizable.strings`, `**/*.xcstrings` |
| Flutter (intl) | `**/l10n/**/*.arb` |
| Flutter (easy_localization) | `**/assets/lang/**/*.json` |
| Rails | `config/locales/**/*.yml` |
| Django | `locale/**/LC_MESSAGES/*.json` |

**Directories always ignored:** `node_modules`, `dist`, `build`, `.git`, `vendor`, `.angular`, `.next`, `Pods`, `DerivedData`, `coverage`.

If your files are in a non-standard location, specify it via `.i18n-fixer.json`:

```json
{
  "preset": "ngx-translate",
  "i18nFilePatterns": ["**/my-custom-path/**/*.json"]
}
```

## Ignoring Lines

Add an `i18n-ignore` comment to suppress false positives on a specific line. Works with any comment style:

```tsx
<Text>Version 1.0.0</Text>  {/* i18n-ignore */}
```
```swift
Text("AppID-XK92")  // i18n-ignore
```
```python
raise Exception("internal-only-error")  # i18n-ignore
```
```dart
Text("debug-mode")  // i18n-ignore
```

## GitHub Actions

```yaml
- uses: zfurkandurum/i18n-fixer@v1
  with:
    preset: react-i18next   # optional, auto-detected by default
    format: console
```

Available inputs: `preset`, `format`, `path`, `no-hardcoded`, `no-missing`, `no-unused`, `args`.

## Pre-commit Hook

Add to your `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/zfurkandurum/i18n-fixer
    rev: v0.2.0
    hooks:
      - id: i18n-fixer
```

## Configuration

Create `.i18n-fixer.json` in your project root to customize behavior:

```json
{
  "preset": "react-i18next",
  "defaultLocale": "en",
  "ignore": ["src/legacy/**"],
  "format": "console",
  "verbose": false
}
```

### Excluding Dynamic Key Namespaces

Some keys are used dynamically at runtime (e.g., error codes returned by the backend). Add them to `unusedKeyIgnorePatterns` to exclude from unused-key reporting:

```json
{
  "preset": "ngx-translate",
  "unusedKeyIgnorePatterns": [
    "ERRORS.*",
    "BACKEND_CODES.*"
  ]
}
```

Supported pattern forms:
- `"ERRORS.*"` — excludes any key starting with `ERRORS.`
- `"*.DEPRECATED"` — excludes any key ending with `.DEPRECATED`
- `"exact.key"` — excludes this exact key

### Adding Project-Specific i18n Patterns

If your project wraps the translate function in a custom method, add extra patterns:

```json
{
  "preset": "ngx-translate",
  "i18nFunctionPatterns": [
    "showError\\(['\"](?P<key>[A-Z][A-Z0-9_.]+)['\"]",
    "showSuccess\\(['\"](?P<key>[A-Z][A-Z0-9_.]+)['\"]",
    "notify\\(['\"](?P<key>[A-Z][A-Z0-9_.]+)['\"]"
  ]
}
```

Or generate a starter config:

```bash
i18n-fixer init
```

## Custom Presets

Create a JSON file with your framework's i18n patterns:

```json
{
  "name": "my-framework",
  "displayName": "My Framework",
  "fileExtensions": [".tsx", ".ts"],
  "i18nFunctionPatterns": [
    "\\bt\\(['\"](?P<key>[^'\"]+)['\"]"
  ],
  "hardcodedStringPatterns": [
    ">[\\s]*(?P<str>[A-Z][a-zA-Z0-9 ,.!?'\\-]{2,})[\\s]*<"
  ],
  "hardcodedStringExclusions": [
    "^https?://", "^[0-9.,]+$"
  ],
  "i18nFilePatterns": ["src/locales/**/*.json"],
  "i18nFileFormat": "json",
  "keyStyle": "nested",
  "keySeparator": ".",
  "projectMarkers": [
    { "file": "package.json", "containsAny": ["my-framework"] }
  ],
  "ignorePatterns": ["**/node_modules/**"]
}
```

Supported `i18nFileFormat` values: `json`, `yaml`, `xml`, `strings`, `arb`, `xcstrings`.

Then use it:

```bash
i18n-fixer --preset ./my-preset.json
```

## AI Prompt Output

The `--format prompt` flag generates a structured Markdown document containing all findings. It automatically saves to `i18n-fix-prompt.md` in your current directory. Paste it into Claude, ChatGPT, or any AI assistant to automatically fix the issues:

```bash
# Auto-saves to i18n-fix-prompt.md
i18n-fixer -f prompt

# Or specify a custom path
i18n-fixer -f prompt -o custom-fix.md
```

The generated prompt includes:
- Missing keys with file locations and target locales
- Unused keys to remove
- Hardcoded strings with suggested i18n keys
- Dynamic keys requiring manual review

## Comparison with Other Tools

### Feature Comparison

| Feature | i18n-fixer | i18next-scanner | eslint-plugin-i18next | i18n-unused | i18n-ally (VS Code) | i18n-tasks (Ruby) |
|---------|-----------|-----------------|----------------------|-------------|-------------------|-------------------|
| Missing keys | :white_check_mark: | :white_check_mark: | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Unused keys | :white_check_mark: | :x: | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Hardcoded strings | :white_check_mark: | :x: | :white_check_mark: | :x: | :white_check_mark: | :x: |
| AI prompt output | :white_check_mark: | :x: | :x: | :x: | :x: | :x: |
| Locale completeness % | :white_check_mark: | :x: | :x: | :x: | :white_check_mark: | :white_check_mark: |
| Duplicate key detection | :white_check_mark: | :x: | :x: | :x: | :x: | :white_check_mark: |
| Key naming lint | :white_check_mark: | :x: | :x: | :x: | :x: | :x: |
| Inline ignore annotation | :white_check_mark: | :x: | :white_check_mark: | :x: | :x: | :x: |
| HTML/template scanning | :white_check_mark: | :x: | :x: | :x: | :white_check_mark: | :x: |
| Auto framework detect | :white_check_mark: | :x: | :x: | :x: | :white_check_mark: | :x: |
| GitHub Actions | :white_check_mark: | :x: | :x: | :x: | :x: | :x: |
| Pre-commit hook | :white_check_mark: | :x: | :white_check_mark: | :x: | :x: | :white_check_mark: |
| Zero config | :white_check_mark: | :x: | :x: | :x: | :white_check_mark: | :x: |
| Single binary | :white_check_mark: | :x: | :x: | :x: | :x: | :x: |
| CI/CD friendly | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x: | :white_check_mark: |

### Framework Support

| Framework | i18n-fixer | i18next-scanner | i18n-unused | i18n-ally |
|-----------|-----------|-----------------|-------------|-----------|
| React (i18next) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| React (intl) | :white_check_mark: | :x: | :x: | :white_check_mark: |
| Vue | :white_check_mark: | :x: | :white_check_mark: | :white_check_mark: |
| Angular | :white_check_mark: | :x: | :x: | :white_check_mark: |
| Svelte | :white_check_mark: | :x: | :x: | :white_check_mark: |
| Next.js | :white_check_mark: | :x: | :x: | :white_check_mark: |
| Flutter | :white_check_mark: | :x: | :x: | :white_check_mark: |
| iOS (Swift) | :white_check_mark: | :x: | :x: | :x: |
| Android | :white_check_mark: | :x: | :x: | :x: |
| React Native | :white_check_mark: | :x: | :x: | :white_check_mark: |
| Django | :white_check_mark: | :x: | :x: | :x: |
| Rails | :white_check_mark: | :x: | :x: | :x: |
| Laravel | :white_check_mark: | :x: | :x: | :x: |
| Go | :white_check_mark: | :x: | :x: | :x: |

### Platform & Distribution

| Aspect | i18n-fixer | i18next-scanner | i18n-unused | i18n-ally |
|--------|-----------|-----------------|-------------|-----------|
| Type | CLI binary | Node.js CLI | Node.js CLI | VS Code extension |
| Runtime | None | Node.js | Node.js | VS Code |
| Install | brew/npx/binary | npm | npm | VS Code Marketplace |
| Language | Go | JavaScript | JavaScript | TypeScript |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

[MIT](LICENSE)
