/*
Copyright © 2025 Tobias Schäfer <tobias.schaefer@orvanta.de>
Licensed under the MIT license.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tschaefer/rpinfo/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long:  "Display version information",
	Run: func(cmd *cobra.Command, args []string) {
		version.Print()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
