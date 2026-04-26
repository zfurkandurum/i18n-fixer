package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zfurkandurum/i18n-fixer/internal/detect"
)

const configTemplate = `{
  "preset": "%s",
  "defaultLocale": "en",
  "ignore": [],
  "format": "console",
  "verbose": false,
  "noHardcoded": false,
  "noMissing": false,
  "noUnused": false,
  "strictUnused": false
}
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a starter .i18n-fixer.json config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := os.Getwd()
		configPath := filepath.Join(dir, ".i18n-fixer.json")

		if _, err := os.Stat(configPath); err == nil {
			return fmt.Errorf(".i18n-fixer.json already exists in %s", dir)
		}

		// Try to auto-detect preset
		presetName := ""
		detected, _ := detect.Detect(dir)
		if len(detected) > 0 {
			presetName = detected[0].Name
			fmt.Printf("Detected framework: %s\n", detected[0].DisplayName)
		}

		content := fmt.Sprintf(configTemplate, presetName)

		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing config file: %w", err)
		}

		fmt.Printf("Created %s\n", configPath)
		if presetName != "" {
			fmt.Printf("Preset set to: %s\n", presetName)
		} else {
			fmt.Println("No framework detected — set \"preset\" manually or use --preset flag")
			fmt.Printf("Available presets: %s\n", strings.Join(presetNames(), ", "))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func presetNames() []string {
	return []string{
		"react-i18next", "react-intl", "vue-i18n",
		"angular", "ngx-translate", "svelte-i18n",
		"next-intl", "nuxt-i18n", "ember-intl",
		"flutter", "ios-swift", "android-kotlin", "react-native",
	}
}
