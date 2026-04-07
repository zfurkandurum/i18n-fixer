package cli

import (
	"fmt"
	"sort"

	"github.com/zfurkandurum/i18n-fixer/internal/preset"
	"github.com/spf13/cobra"
)

var presetsCmd = &cobra.Command{
	Use:   "presets",
	Short: "List available built-in framework presets",
	Run: func(cmd *cobra.Command, args []string) {
		all := preset.All()

		names := make([]string, 0, len(all))
		for name := range all {
			names = append(names, name)
		}
		sort.Strings(names)

		fmt.Println("Available presets:")
		fmt.Println()
		for _, name := range names {
			p := all[name]
			fmt.Printf("  %-20s %s\n", name, p.DisplayName)
		}
		fmt.Println()
		fmt.Println("Use: i18n-fixer --preset <name>")
		fmt.Println("Or create a custom preset JSON file: i18n-fixer --preset ./my-preset.json")
	},
}

func init() {
	rootCmd.AddCommand(presetsCmd)
}
