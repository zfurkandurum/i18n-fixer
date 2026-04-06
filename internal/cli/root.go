package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionStr string
	commitStr  string
	dateStr    string
)

func SetVersionInfo(version, commit, date string) {
	versionStr = version
	commitStr = commit
	dateStr = date
}

var rootCmd = &cobra.Command{
	Use:   "i18n-fixer [flags] [path]",
	Short: "Find and fix i18n issues across any framework",
	Long: `i18n-fixer is a framework-agnostic CLI tool that finds hardcoded strings,
missing translation keys, and unused translations across React, Vue, Angular,
Flutter, iOS, Android, and more.

It auto-detects your framework and generates AI-ready fix prompts.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("i18n-fixer %s\ncommit: %s\nbuilt at: %s\n", versionStr, commitStr, dateStr)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Global flags
	rootCmd.PersistentFlags().StringP("preset", "p", "", "framework preset name or path to custom preset JSON")
	rootCmd.PersistentFlags().StringP("format", "f", "console", "output format: console, json, prompt")
	rootCmd.PersistentFlags().StringP("output", "o", "", "write report to file")
	rootCmd.PersistentFlags().String("default-locale", "", "only check missing keys against this locale")
	rootCmd.PersistentFlags().Bool("no-hardcoded", false, "skip hardcoded string detection")
	rootCmd.PersistentFlags().Bool("no-missing", false, "skip missing key detection")
	rootCmd.PersistentFlags().Bool("no-unused", false, "skip unused key detection")
	rootCmd.PersistentFlags().StringSlice("ignore", nil, "additional glob patterns to ignore")
	rootCmd.PersistentFlags().Bool("strict-unused", false, "disable dynamic key heuristic exclusion")
	rootCmd.PersistentFlags().Bool("verbose", false, "show detailed scanning progress")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")
}

func Execute() error {
	return rootCmd.Execute()
}
