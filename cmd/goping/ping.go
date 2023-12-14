package goping

import (
	"fmt"
	"time"

	"github.com/celestiaorg/goping/client"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	flagLogLevel       = "log-level"
	flagProductionMode = "production-mode"

	flagPacketsCount  = "packets-count"
	flagPacketTimeout = "packet-timeout"
	flagMode          = "mode"
	flagQuiet         = "quiet"
)

const (
	pingModePacketloss = "packetloss"
	pingModeLatency    = "latency"
	pingModeJitter     = "jitter"
)

var flagsPing struct {
	logLevel       string
	productionMode bool

	packetsCount  int
	packetTimeout time.Duration
	mode          string
	quiet         bool
}

func init() {
	rootCmd.AddCommand(pingCmd)

	pingCmd.PersistentFlags().IntVarP(&flagsPing.packetsCount, flagPacketsCount, "c", 4, "number of packets to send")
	pingCmd.PersistentFlags().DurationVarP(&flagsPing.packetTimeout, flagPacketTimeout, "t", 10*time.Millisecond, "timeout for each packet")
	pingCmd.PersistentFlags().BoolVarP(&flagsPing.quiet, flagQuiet, "q", false, "quiet mode (e.g. only print the summary)")
	pingCmd.PersistentFlags().StringVarP(&flagsPing.mode, flagMode, "m", pingModePacketloss, "mode (e.g. packetloss, latency, jitter)")

	pingCmd.PersistentFlags().StringVar(&flagsPing.logLevel, flagLogLevel, zap.InfoLevel.String(), "log level (e.g. debug, info, warn, error, dpanic, panic, fatal)")
	pingCmd.PersistentFlags().BoolVar(&flagsPing.productionMode, flagProductionMode, false, "production mode (e.g. disable debug logs)")
}

var pingCmd = &cobra.Command{
	Use:   "ping <server-address>",
	Short: "ping the goping server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := getLogger(flagsPing.logLevel, flagsPing.productionMode)
		if err != nil {
			return err
		}
		defer func() {
			// The error is ignored because of this issue: https://github.com/uber-go/zap/issues/328
			_ = logger.Sync()
		}()

		if flagsPing.quiet {
			logger = zap.NewNop()
		}

		serverAddr := args[0]
		logger.Info("pinging server", zap.String("address", serverAddr))

		if flagsPing.mode == pingModePacketloss {

			packetLossPercentage, err := client.MeasurePacketloss(serverAddr, flagsPing.packetsCount, flagsPing.packetTimeout, logger)
			if err != nil {
				return err
			}

			logger.Info("Packet loss percentage", zap.Float64("percentage", packetLossPercentage))

			if flagsPing.quiet {
				fmt.Print(packetLossPercentage)
			}
			return nil
		}

		if flagsPing.mode == pingModeLatency {
			latency, err := client.MeasureLatency(serverAddr, flagsPing.packetsCount, flagsPing.packetTimeout, logger)
			if err != nil {
				return err
			}

			logger.Info("Latency", zap.Duration("latency", latency))

			if flagsPing.quiet {
				fmt.Print(latency.String())
			}
			return nil
		}

		if flagsPing.mode == pingModeJitter {
			jitter, err := client.MeasureJitter(serverAddr, flagsPing.packetsCount, flagsPing.packetTimeout, logger)
			if err != nil {
				return err
			}

			logger.Info("Jitter", zap.Duration("jitter", jitter))

			if flagsPing.quiet {
				fmt.Print(jitter)
			}
		}

		return nil
	},
}
