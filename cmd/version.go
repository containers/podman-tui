package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	appVersion = "0.4.0"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Display %s version and exit.\n", appName),
	Long:  fmt.Sprintf("Display %s version and exit.\n", appName),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v%s\n", appName, appVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
