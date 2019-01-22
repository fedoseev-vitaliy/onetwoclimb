package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var config Config
var l = logrus.New()

func init() {
	RootCmd.AddCommand(CMD)
	CMD.Flags().AddFlagSet(config.Flags())
}

var CMD = &cobra.Command{
	Use:   "server",
	Short: "run server",
	RunE: func(cmd *cobra.Command, args []string) error {
		l.Info("start 12Climb server")
		defer l.Info("stop 12Climb server")

		return nil
	},
}
