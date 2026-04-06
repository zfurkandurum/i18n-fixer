package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/i18n-fixer/i18n-fixer/internal/types"
)

// PromptReporter generates an AI-ready fix prompt in Markdown.
type PromptReporter struct{}

func (r *PromptReporter) Report(result *types.AuditResult, w io.Writer) error {
	fmt.Fprintf(w, "# i18n Fix Request\n\n")
	fmt.Fprintf(w, "Project: `%s`\n", result.Metadata.RootDir)
	fmt.Fprintf(w, "Framework: %s | Locales: %s\n", result.Metadata.Preset, strings.Join(result.Summary.Locales, ", "))
	fmt.Fprintf(w, "Translation files: %s format\n\n", result.Metadata.I18nFileFormat)

	// Missing Keys
	if len(result.MissingKeys) > 0 {
		fmt.Fprintf(w, "## 1. Missing Translation Keys (%d)\n\n", len(result.MissingKeys))
		fmt.Fprintf(w, "Keys used in source code but missing from translation files:\n\n")
		fmt.Fprintf(w, "| # | Key | Used In | Missing From |\n")
		fmt.Fprintf(w, "|---|-----|---------|-------------|\n")
		for i, issue := range result.MissingKeys {
			locations := make([]string, 0, len(issue.UsedIn))
			for _, loc := range issue.UsedIn {
				locations = append(locations, fmt.Sprintf("`%s:%d`", loc.File, loc.Line))
			}
			fmt.Fprintf(w, "| %d | `%s` | %s | %s |\n",
				i+1, issue.Key, strings.Join(locations, ", "), strings.Join(issue.MissingFromLocales, ", "))
		}
		fmt.Fprintf(w, "\n**Action**: Add each missing key to the listed locale files with appropriate translations.\n\n")
	}

	// Unused Keys
	if len(result.UnusedKeys) > 0 {
		fmt.Fprintf(w, "## 2. Unused Translation Keys (%d)\n\n", len(result.UnusedKeys))
		fmt.Fprintf(w, "Keys in translation files but never referenced in source code:\n\n")
		fmt.Fprintf(w, "| # | Key | Defined In |\n")
		fmt.Fprintf(w, "|---|-----|------------|\n")
		for i, issue := range result.UnusedKeys {
			files := make([]string, 0, len(issue.DefinedIn))
			for _, loc := range issue.DefinedIn {
				files = append(files, fmt.Sprintf("`%s`", loc.File))
			}
			fmt.Fprintf(w, "| %d | `%s` | %s |\n", i+1, issue.Key, strings.Join(files, ", "))
		}
		fmt.Fprintf(w, "\n**Action**: Remove these keys from all locale files, or verify they are used via dynamic keys.\n\n")
	}

	// Hardcoded Strings
	if len(result.HardcodedStrings) > 0 {
		fmt.Fprintf(w, "## 3. Hardcoded Strings (%d)\n\n", len(result.HardcodedStrings))
		fmt.Fprintf(w, "User-facing strings not wrapped in i18n functions:\n\n")
		fmt.Fprintf(w, "| # | String | File:Line | Suggested Key |\n")
		fmt.Fprintf(w, "|---|--------|-----------|---------------|\n")
		for i, issue := range result.HardcodedStrings {
			loc := issue.Occurrences[0]
			fmt.Fprintf(w, "| %d | \"%s\" | `%s:%d` | `%s` |\n",
				i+1, issue.Value, loc.File, loc.Line, issue.SuggestedKey)
		}
		fmt.Fprintf(w, "\n**Action**: For each hardcoded string:\n")
		fmt.Fprintf(w, "1. Add the suggested key to all locale translation files.\n")
		fmt.Fprintf(w, "2. Replace the hardcoded string in source with the i18n function call.\n\n")
	}

	// Dynamic Keys
	if len(result.DynamicKeys) > 0 {
		fmt.Fprintf(w, "## 4. Dynamic Keys — Manual Review (%d)\n\n", len(result.DynamicKeys))
		fmt.Fprintf(w, "These expressions use computed keys that cannot be statically analyzed:\n\n")
		fmt.Fprintf(w, "| # | Expression | File:Line |\n")
		fmt.Fprintf(w, "|---|-----------|----------|\n")
		for i, dk := range result.DynamicKeys {
			fmt.Fprintf(w, "| %d | `%s` | `%s:%d` |\n", i+1, dk.RawExpression, dk.File, dk.Line)
		}
		fmt.Fprintf(w, "\n**Action**: Verify all possible runtime key values exist in all locale files.\n\n")
	}

	// Locale Completeness
	if len(result.LocaleCoverage) > 0 {
		sectionNum := 5
		fmt.Fprintf(w, "## %d. Locale Completeness\n\n", sectionNum)
		fmt.Fprintf(w, "| Locale | Keys | Total | Coverage |\n")
		fmt.Fprintf(w, "|--------|------|-------|----------|\n")
		for _, c := range result.LocaleCoverage {
			fmt.Fprintf(w, "| %s | %d | %d | %.1f%% |\n", c.Locale, c.HasKeys, c.TotalKeys, c.Percentage)
		}
		fmt.Fprintf(w, "\n**Overall completeness: %.1f%%**\n\n", result.Summary.OverallCompleteness)
	}

	// Duplicate Keys
	if len(result.DuplicateKeys) > 0 {
		sectionNum := 6
		fmt.Fprintf(w, "## %d. Duplicate Keys (%d)\n\n", sectionNum, len(result.DuplicateKeys))
		fmt.Fprintf(w, "Same key with conflicting values in the same locale:\n\n")
		for i, issue := range result.DuplicateKeys {
			fmt.Fprintf(w, "**%d. `%s` [%s]**\n", i+1, issue.Key, issue.Locale)
			for _, v := range issue.Values {
				fmt.Fprintf(w, "- \"%s\" in `%s`\n", v.Value, v.File)
			}
			fmt.Fprintf(w, "\n")
		}
		fmt.Fprintf(w, "**Action**: Resolve conflicting values — keep one, remove duplicates.\n\n")
	}

	// Key Naming Issues
	if len(result.KeyNamingIssues) > 0 {
		sectionNum := 7
		fmt.Fprintf(w, "## %d. Key Naming Violations (%d)\n\n", sectionNum, len(result.KeyNamingIssues))
		fmt.Fprintf(w, "Expected convention: **%s**\n\n", result.KeyNamingIssues[0].Expected)
		fmt.Fprintf(w, "| # | Key | File |\n")
		fmt.Fprintf(w, "|---|-----|------|\n")
		for i, issue := range result.KeyNamingIssues {
			fmt.Fprintf(w, "| %d | `%s` | `%s` |\n", i+1, issue.Key, issue.File)
		}
		fmt.Fprintf(w, "\n**Action**: Rename keys to follow the %s convention.\n\n", result.KeyNamingIssues[0].Expected)
	}

	fmt.Fprintf(w, "---\n")
	fmt.Fprintf(w, "*Generated by i18n-fixer %s on %s*\n", result.Metadata.Version, result.Metadata.Timestamp)

	return nil
}
