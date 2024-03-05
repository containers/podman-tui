package cmd

import (
	"fmt"
	"io"
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

func run(cmd *cobra.Command, args []string) error { //nolint:cyclop,revive
	var (
		logOutput = io.Discard
		runLog    = fmt.Sprintf("starting %s version %s", appName, appVersion)
	)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Default level is info
	debugLevel, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}

	if debugLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		runLog += " in debug mode"

		// init logger
		logfile, err := cmd.Flags().GetString("log-file")
		if err != nil {
			return err
		}

		logFD, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		if err != nil {
			return err
		}

		defer logFD.Close()

		logOutput = logFD
	}

	logrus.SetOutput(logOutput)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: logOutput, TimeFormat: time.RFC3339})

	log.Info().Msg(runLog)
	// check if CONTAINER_PASSPHRASE environment variable is set and not empty
	// otherwise set with value dummy value
	// its required since podman/pkg/podman is using terminal package
	// that is writing directly to os.Stdout and reading from os.Stdin
	// for Phassphrase
	setSSHIdentityPassphrase := true

	v, found := os.LookupEnv("CONTAINER_PASSPHRASE")
	if found {
		if v != "" {
			setSSHIdentityPassphrase = false
		}
	}

	if setSSHIdentityPassphrase {
		emptyPassphrase := "__empty__"

		log.Debug().Msgf("env set CONTAINER_PASSPHRASE=%q", emptyPassphrase)

		err := os.Setenv("CONTAINER_PASSPHRASE", emptyPassphrase)
		if err != nil {
			return err
		}
	}

	app := app.NewApp(appName, appVersion)
	if err := app.Run(); err != nil {
		if setSSHIdentityPassphrase {
			os.Unsetenv("CONTAINER_PASSPHRASE")
		}

		return err
	}

	// unset CONTAINER_PASSPHRASE environment variable if we have set it
	// after application exits
	if setSSHIdentityPassphrase {
		os.Unsetenv("CONTAINER_PASSPHRASE")
	}

	return nil
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	defaultLogFile := appName + ".log"

	rootCmd.Flags().BoolP("debug", "d", false, "Run application in debug mode")
	rootCmd.Flags().StringP("log-file", "l", defaultLogFile, "Application runtime log file")
}
