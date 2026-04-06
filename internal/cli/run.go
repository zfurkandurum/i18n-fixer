package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [path]",
	Short: "Scan project for i18n issues (default command)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		// TODO: wire up full pipeline in Phase 7
		fmt.Printf("Scanning %s for i18n issues...\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Make "run" the default command when no subcommand is given
	rootCmd.RunE = runCmd.RunE
}
