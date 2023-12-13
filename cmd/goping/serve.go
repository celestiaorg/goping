package goping

import (
	"os"
	"os/signal"

	"github.com/celestiaorg/goping/server"
	"github.com/spf13/cobra"
)

const (
	flagServeAddr = "serve-addr"
)

var flagsServe struct {
	serveAddr      string
	originAllowed  string
	logLevel       string
	productionMode bool
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringVar(&flagsServe.serveAddr, flagServeAddr, ":8080", "address to serve on")
	serveCmd.PersistentFlags().StringVar(&flagsServe.logLevel, flagLogLevel, "info", "log level (e.g. debug, info, warn, error, dpanic, panic, fatal)")
	serveCmd.PersistentFlags().BoolVar(&flagsServe.productionMode, flagProductionMode, false, "production mode (e.g. disable debug logs)")
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the goping server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := getLogger(flagsServe.logLevel, flagsServe.productionMode)
		if err != nil {
			return err
		}
		defer func() {
			// The error is ignored because of this issue: https://github.com/uber-go/zap/issues/328
			_ = logger.Sync()
		}()

		logger.Info("Starting the API server...")

		cancel, err := server.Serve(flagsServe.serveAddr, logger)
		if err != nil {
			return err
		}
		defer cancel()

		// Handle interrupt signal (Ctrl+C)
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		<-signalChan

		return nil
	},
}
