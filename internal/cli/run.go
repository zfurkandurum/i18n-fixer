package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/i18n-fixer/i18n-fixer/internal/analyzer"
	"github.com/i18n-fixer/i18n-fixer/internal/config"
	"github.com/i18n-fixer/i18n-fixer/internal/detect"
	"github.com/i18n-fixer/i18n-fixer/internal/parser"
	"github.com/i18n-fixer/i18n-fixer/internal/preset"
	"github.com/i18n-fixer/i18n-fixer/internal/reporter"
	"github.com/i18n-fixer/i18n-fixer/internal/scanner"
	"github.com/i18n-fixer/i18n-fixer/internal/types"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [path]",
	Short: "Scan project for i18n issues (default command)",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runAudit,
}

func init() {
	rootCmd.AddCommand(runCmd)
	// Default behavior: run audit when no subcommand given
	rootCmd.Args = cobra.MaximumNArgs(1)
	rootCmd.RunE = runAudit
}

func runAudit(cmd *cobra.Command, args []string) error {
	startTime := time.Now()

	// Resolve project path
	rootDir := "."
	if len(args) > 0 {
		rootDir = args[0]
	}
	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	// Load config
	cfg, err := config.Load(rootDir)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// CLI flags override config
	verbose, _ := cmd.Flags().GetBool("verbose")
	format, _ := cmd.Flags().GetString("format")
	output, _ := cmd.Flags().GetString("output")
	presetFlag, _ := cmd.Flags().GetString("preset")
	noHardcoded, _ := cmd.Flags().GetBool("no-hardcoded")
	noMissing, _ := cmd.Flags().GetBool("no-missing")
	noUnused, _ := cmd.Flags().GetBool("no-unused")

	if format == "" {
		format = cfg.Format
	}
	if format == "" {
		format = "console"
	}
	if presetFlag == "" {
		presetFlag = cfg.Preset
	}

	// Detect or load preset
	var presets []types.FrameworkPreset
	if presetFlag != "" {
		// Try as built-in preset name first
		p, err := preset.Get(presetFlag)
		if err != nil {
			// Try as custom file path
			p, err = preset.LoadCustom(presetFlag)
			if err != nil {
				return fmt.Errorf("loading preset %q: %w", presetFlag, err)
			}
		}
		presets = []types.FrameworkPreset{p}
	} else {
		// Auto-detect
		presets, err = detect.Detect(rootDir)
		if err != nil {
			return fmt.Errorf("detecting framework: %w", err)
		}
		if len(presets) == 0 {
			return fmt.Errorf("no framework detected in %s\n\nUse --preset <name> or create .i18n-fixer.json\nRun 'i18n-fixer presets' to see available presets", rootDir)
		}
	}

	if verbose {
		names := make([]string, 0, len(presets))
		for _, p := range presets {
			names = append(names, p.DisplayName)
		}
		fmt.Fprintf(os.Stderr, "Detected: %s\n", strings.Join(names, ", "))
	}

	// Run audit for each detected preset
	var allResults []*types.AuditResult

	for _, p := range presets {
		if verbose {
			fmt.Fprintf(os.Stderr, "Scanning with preset: %s\n", p.DisplayName)
		}

		// Scan source files
		scanResult, err := scanner.Scan(rootDir, p)
		if err != nil {
			return fmt.Errorf("scanning source files: %w", err)
		}

		// Parse i18n files
		i18nEntries, i18nFileCount, err := parseI18nFiles(rootDir, p)
		if err != nil {
			return fmt.Errorf("parsing i18n files: %w", err)
		}

		// Analyze
		result := analyzer.Analyze(scanResult, i18nEntries, p.KeySeparator, analyzer.Options{
			NoMissing:   noMissing || cfg.NoMissing,
			NoUnused:    noUnused || cfg.NoUnused,
			NoHardcoded: noHardcoded || cfg.NoHardcoded,
		})

		result.Summary.FilesScanned = countSourceFiles(rootDir, p)
		result.Summary.I18nFilesLoaded = i18nFileCount
		result.Metadata = types.AuditMetadata{
			Tool:           "i18n-fixer",
			Version:        versionStr,
			Timestamp:      time.Now().UTC().Format(time.RFC3339),
			Preset:         p.DisplayName,
			RootDir:        rootDir,
			I18nFileFormat: p.I18nFileFormat,
			Duration:       time.Since(startTime).Milliseconds(),
		}

		allResults = append(allResults, result)
	}

	// Merge results if multiple presets
	finalResult := mergeResults(allResults)

	// Report
	rep, err := reporter.New(format)
	if err != nil {
		return err
	}

	// Auto-generate output file for prompt format (AI prompt is meant to be
	// copied to an AI agent, not read in the terminal)
	if format == "prompt" && output == "" {
		output = "i18n-fix-prompt.md"
	}

	var w io.Writer = os.Stdout
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
		fmt.Fprintf(os.Stderr, "Report written to %s\n", output)
	}

	if err := rep.Report(finalResult, w); err != nil {
		return fmt.Errorf("generating report: %w", err)
	}

	// Return special error for CI: issues found but not a failure
	total := finalResult.Summary.MissingKeyCount + finalResult.Summary.UnusedKeyCount +
		finalResult.Summary.HardcodedStringCount
	if total > 0 {
		return &IssuesFoundError{Count: total}
	}

	return nil
}

func parseI18nFiles(rootDir string, p types.FrameworkPreset) ([]types.I18nEntry, int, error) {
	var allEntries []types.I18nEntry
	fileCount := 0

	for _, pattern := range p.I18nFilePatterns {
		fullPattern := filepath.Join(rootDir, pattern)
		matches, err := doublestar.FilepathGlob(fullPattern)
		if err != nil {
			continue
		}

		for _, match := range matches {
			entries, err := parser.Parse(match, p.I18nFileFormat, p.KeySeparator)
			if err != nil {
				continue
			}
			allEntries = append(allEntries, entries...)
			fileCount++
		}
	}

	return allEntries, fileCount, nil
}

func countSourceFiles(rootDir string, p types.FrameworkPreset) int {
	count := 0
	extSet := make(map[string]bool)
	for _, ext := range p.FileExtensions {
		extSet[ext] = true
	}
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if extSet[filepath.Ext(path)] {
			count++
		}
		return nil
	})
	return count
}

func mergeResults(results []*types.AuditResult) *types.AuditResult {
	if len(results) == 1 {
		return results[0]
	}

	merged := &types.AuditResult{}
	for _, r := range results {
		merged.MissingKeys = append(merged.MissingKeys, r.MissingKeys...)
		merged.UnusedKeys = append(merged.UnusedKeys, r.UnusedKeys...)
		merged.HardcodedStrings = append(merged.HardcodedStrings, r.HardcodedStrings...)
		merged.DynamicKeys = append(merged.DynamicKeys, r.DynamicKeys...)
		merged.Summary.FilesScanned += r.Summary.FilesScanned
		merged.Summary.I18nFilesLoaded += r.Summary.I18nFilesLoaded
		merged.Summary.TotalKeys += r.Summary.TotalKeys
		merged.Summary.MissingKeyCount += r.Summary.MissingKeyCount
		merged.Summary.UnusedKeyCount += r.Summary.UnusedKeyCount
		merged.Summary.HardcodedStringCount += r.Summary.HardcodedStringCount
		merged.Summary.DynamicKeyCount += r.Summary.DynamicKeyCount
		for _, l := range r.Summary.Locales {
			merged.Summary.Locales = appendUnique(merged.Summary.Locales, l)
		}
		merged.Metadata = r.Metadata // use last preset's metadata
	}
	return merged
}

// IssuesFoundError indicates the scan completed successfully but found issues.
type IssuesFoundError struct {
	Count int
}

func (e *IssuesFoundError) Error() string {
	return fmt.Sprintf("found %d i18n issues", e.Count)
}

func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}
