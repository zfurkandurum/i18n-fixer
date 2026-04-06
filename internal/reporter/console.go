package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/i18n-fixer/i18n-fixer/internal/types"
)

// ConsoleReporter writes a human-readable colored report to the terminal.
type ConsoleReporter struct{}

func (r *ConsoleReporter) Report(result *types.AuditResult, w io.Writer) error {
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "══════════════════════════════════════════════════════\n")
	fmt.Fprintf(w, "  i18n-fixer — Internationalization Audit Report\n")
	fmt.Fprintf(w, "  Preset: %s | Locales: %s\n", result.Metadata.Preset, strings.Join(result.Summary.Locales, ", "))
	fmt.Fprintf(w, "══════════════════════════════════════════════════════\n\n")

	fmt.Fprintf(w, "Scanned %d source files, loaded %d i18n files (%d locales)\n",
		result.Summary.FilesScanned, result.Summary.I18nFilesLoaded, len(result.Summary.Locales))
	fmt.Fprintf(w, "Found %d i18n keys defined\n\n", result.Summary.TotalKeys)

	// Missing Keys
	if len(result.MissingKeys) > 0 {
		fmt.Fprintf(w, "━━━ MISSING KEYS (%d) ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n", len(result.MissingKeys))
		for _, issue := range result.MissingKeys {
			fmt.Fprintf(w, "  %s\n", issue.Key)
			fmt.Fprintf(w, "    Used in:\n")
			for _, loc := range issue.UsedIn {
				fmt.Fprintf(w, "      %s:%d\n", loc.File, loc.Line)
			}
			fmt.Fprintf(w, "    Missing from: %s\n\n", strings.Join(issue.MissingFromLocales, ", "))
		}
	}

	// Unused Keys
	if len(result.UnusedKeys) > 0 {
		fmt.Fprintf(w, "━━━ UNUSED KEYS (%d) ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n", len(result.UnusedKeys))
		for _, issue := range result.UnusedKeys {
			files := make([]string, 0, len(issue.DefinedIn))
			for _, loc := range issue.DefinedIn {
				files = append(files, loc.File)
			}
			fmt.Fprintf(w, "  %-40s %s\n", issue.Key, strings.Join(files, ", "))
		}
		fmt.Fprintf(w, "\n")
	}

	// Hardcoded Strings
	if len(result.HardcodedStrings) > 0 {
		fmt.Fprintf(w, "━━━ HARDCODED STRINGS (%d) ━━━━━━━━━━━━━━━━━━━━━━━━━\n\n", len(result.HardcodedStrings))
		for _, issue := range result.HardcodedStrings {
			fmt.Fprintf(w, "  \"%s\"\n", issue.Value)
			for _, occ := range issue.Occurrences {
				fmt.Fprintf(w, "    %s:%d\n", occ.File, occ.Line)
			}
			fmt.Fprintf(w, "    Suggested key: %s\n\n", issue.SuggestedKey)
		}
	}

	// Dynamic Keys
	if len(result.DynamicKeys) > 0 {
		fmt.Fprintf(w, "━━━ DYNAMIC KEYS (%d) — manual review needed ━━━━━━━\n\n", len(result.DynamicKeys))
		for _, dk := range result.DynamicKeys {
			fmt.Fprintf(w, "  %s\n", dk.RawExpression)
			fmt.Fprintf(w, "    %s:%d\n\n", dk.File, dk.Line)
		}
	}

	// Summary
	total := result.Summary.MissingKeyCount + result.Summary.UnusedKeyCount +
		result.Summary.HardcodedStringCount + result.Summary.DynamicKeyCount

	fmt.Fprintf(w, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Fprintf(w, "  Summary:\n")
	fmt.Fprintf(w, "  ┌──────────────────────┬───────┐\n")
	fmt.Fprintf(w, "  │ Missing keys         │ %5d │\n", result.Summary.MissingKeyCount)
	fmt.Fprintf(w, "  │ Unused keys          │ %5d │\n", result.Summary.UnusedKeyCount)
	fmt.Fprintf(w, "  │ Hardcoded strings    │ %5d │\n", result.Summary.HardcodedStringCount)
	fmt.Fprintf(w, "  │ Dynamic keys (warn)  │ %5d │\n", result.Summary.DynamicKeyCount)
	fmt.Fprintf(w, "  │ Total issues         │ %5d │\n", total)
	fmt.Fprintf(w, "  └──────────────────────┴───────┘\n")
	fmt.Fprintf(w, "\n")

	return nil
}
