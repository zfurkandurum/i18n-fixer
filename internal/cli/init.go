package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a starter .i18n-fixer.json config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implement config generation in Phase 2
		fmt.Println("Generating .i18n-fixer.json...")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
