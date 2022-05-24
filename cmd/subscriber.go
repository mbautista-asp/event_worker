package cmd

import (
	"github.com/charliemcelfresh/event_worker/internal/subscriber"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Run event subscriber",
	Run: func(cmd *cobra.Command, args []string) {
		subscriber.Subscribe()
	},
}
