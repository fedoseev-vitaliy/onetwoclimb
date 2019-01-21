package cmd

import "github.com/spf13/cobra"

var config Config

func init() {
	RootCmd.AddCommand(CMD)
	CMD.Flags().AddFlagSet(config.Flags())
}

var CMD = &cobra.Command{
	Use:   "server",
	Short: "run server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
