package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
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
