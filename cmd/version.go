package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	appVersion = "0.17.0-dev"
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Display %s version and exit.\n", appName),
	Long:  fmt.Sprintf("Display %s version and exit.\n", appName),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v%s\n", appName, appVersion) //nolint:forbidigo
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
