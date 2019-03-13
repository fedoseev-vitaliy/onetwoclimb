package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"

	"github.com/onetwoclimb/cmd/migration"
	"github.com/onetwoclimb/cmd/server"
)

var RootCmd = &cobra.Command{
	Use:   "OneTwoClimbAPI",
	Short: "OneTwoClimb swagger API",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("something goes wrong")
		return
	}
}

func init() {
	RootCmd.AddCommand(server.Cmd)
	RootCmd.AddCommand(migration.Migration)
}
