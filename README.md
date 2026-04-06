# i18n-fixer

[![CI](https://github.com/i18n-fixer/i18n-fixer/actions/workflows/ci.yml/badge.svg)](https://github.com/i18n-fixer/i18n-fixer/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/i18n-fixer/i18n-fixer.svg)](https://pkg.go.dev/github.com/i18n-fixer/i18n-fixer)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> Framework-agnostic CLI tool that finds hardcoded strings, missing i18n keys, and unused translations. Generates AI-ready fix prompts. Works with every major frontend and mobile framework.

## Supported Frameworks

| Category | Frameworks |
|----------|-----------|
| **React** | react-i18next, react-intl (FormatJS) |
| **Vue** | vue-i18n (Vue 2 & 3) |
| **Angular** | @angular/localize, @ngx-translate |
| **Svelte** | svelte-i18n |
| **Next.js** | next-intl |
| **Nuxt** | @nuxtjs/i18n |
| **Ember** | ember-intl |
| **Flutter** | intl / ARB files |
| **iOS** | NSLocalizedString, String(localized:) |
| **Android** | getString(R.string.x), XML resources |
| **React Native** | i18next, react-intl |
| **Custom** | Any framework via custom preset JSON |

## Features

- **Auto-detects** your framework — zero config needed
- **13 built-in presets** covering all major frontend & mobile frameworks
- Finds **missing** translation keys (used in code, absent from i18n files)
- Finds **unused** translation keys (in i18n files, never referenced in code)
- Detects **hardcoded** user-facing strings not wrapped in i18n functions
- **Locale completeness %** — see translation coverage per locale at a glance
- **Duplicate key detection** — finds conflicting values for the same key
- **Key naming lint** — enforces UPPER_SNAKE, lower.dot, camelCase, or kebab-case
- Generates **AI-ready fix prompts** — paste into Claude/GPT to auto-fix
- **HTML template scanning** — catches keys in Angular pipes, Vue templates, etc.
- **Fast** — parallel scanning via goroutines
- **Single binary** — no Node.js, Python, or any runtime needed
- **Extensible** — add any framework via custom preset JSON
- **CI/CD ready** — exit codes + JSON output for pipelines

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
go install github.com/i18n-fixer/i18n-fixer/cmd/i18n-fixer@latest
```

### Binary Download

Download from [GitHub Releases](https://github.com/i18n-fixer/i18n-fixer/releases).

## Usage

```
i18n-fixer [flags] [path]

Commands:
  run          Scan project for i18n issues (default)
  init         Generate a starter .i18n-fixer.json config
  presets      List available built-in presets
  version      Print version information

Flags:
  -p, --preset <name>     Framework preset (auto-detected by default)
  -f, --format <type>     Output: console, json, prompt (default: console)
  -o, --output <path>     Write report to file
      --no-hardcoded      Skip hardcoded string detection
      --no-missing        Skip missing key detection
      --no-unused         Skip unused key detection
      --verbose           Show scanning progress
      --no-color          Disable colored output
```

### Examples

```bash
# Auto-detect and scan current directory
i18n-fixer

# Scan a specific directory
i18n-fixer ./frontend

# Use a specific preset
i18n-fixer --preset vue-i18n

# Generate AI prompt (auto-saves to i18n-fix-prompt.md)
i18n-fixer -f prompt

# Or specify a custom output path
i18n-fixer -f prompt -o my-fix.md

# JSON output for CI pipelines
i18n-fixer -f json -o report.json

# Only check for missing keys
i18n-fixer --no-hardcoded --no-unused

# Use custom preset for unsupported framework
i18n-fixer --preset ./my-preset.json
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
| HTML template scanning | :white_check_mark: | :x: | :x: | :x: | :white_check_mark: | :x: |
| Auto framework detect | :white_check_mark: | :x: | :x: | :x: | :white_check_mark: | :x: |
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
