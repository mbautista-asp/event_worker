package cmd

import (
	"github.com/charliemcelfresh/event_worker/internal/config"
	"github.com/charliemcelfresh/event_worker/internal/subscriber"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Run event subscriber",
	Run: func(cmd *cobra.Command, args []string) {
		store := subscriber.NewStore(config.DBPool)
		s := subscriber.New(store, time.Second*5, 10)
		s.Subscribe()
	},
}
