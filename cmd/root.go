package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/containers/podman-tui/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

const (
	appName = "podman-tui"
)

var rootCmd = &cobra.Command{
	Use:   appName,
	Short: "Podman terminal user interface",
	Long:  `Podman terminal user interface`,
	RunE:  run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func run(cmd *cobra.Command, args []string) error {
	runLog := fmt.Sprintf("starting %s version %s", appName, appVersion)
	// init logger
	logfile, err := cmd.Flags().GetString("log-file")
	if err != nil {
		return nil
	}
	logFD, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer logFD.Close()

	logrus.SetOutput(logFD)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: logFD, TimeFormat: time.RFC3339})

	// Default level is info
	debugLevel, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return nil
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debugLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		runLog = runLog + " in debug mode"
	}

	log.Info().Msg(runLog)
	app := app.NewApp()
	if err := app.Run(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	defaultLogFile := appName + ".log"
	rootCmd.Flags().BoolP("debug", "d", false, "Run application in debug mode")
	rootCmd.Flags().StringP("log-file", "l", defaultLogFile, "Application runtime log file")
}
