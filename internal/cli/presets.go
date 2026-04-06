package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var presetsCmd = &cobra.Command{
	Use:   "presets",
	Short: "List available built-in framework presets",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: list from registry in Phase 2
		fmt.Println("Available presets:")
		fmt.Println("  react-i18next    React + i18next")
		fmt.Println("  react-intl       React + FormatJS")
		fmt.Println("  vue-i18n         Vue 2/3")
		fmt.Println("  angular          @angular/localize")
		fmt.Println("  ngx-translate    Angular + ngx-translate")
		fmt.Println("  svelte-i18n      Svelte")
		fmt.Println("  next-intl        Next.js")
		fmt.Println("  nuxt-i18n        Nuxt.js")
		fmt.Println("  ember-intl       Ember.js")
		fmt.Println("  flutter          Flutter/Dart")
		fmt.Println("  ios-swift        iOS Swift")
		fmt.Println("  android-kotlin   Android Kotlin/Java")
		fmt.Println("  react-native     React Native")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(presetsCmd)
}
