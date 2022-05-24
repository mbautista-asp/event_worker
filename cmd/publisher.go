package cmd

import (
	"strconv"

	"github.com/charliemcelfresh/event_worker/internal/config"

	"github.com/charliemcelfresh/event_worker/internal/publisher"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(publisherCmd)
}

var publisherCmd = &cobra.Command{
	Use:   "publish",
	Short: "Run event publisher",
	Run: func(cmd *cobra.Command, args []string) {
		countOfMsgs, err := strconv.Atoi(args[0])
		if err != nil {
			config.Logger.Panic(err)
		}
		publisher.Publish(countOfMsgs)
	},
}
