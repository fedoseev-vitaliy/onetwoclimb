package cmd

import (
	"flag"

	"github.com/go-openapi/loads"
	"github.com/onetwoclimb/internal/server/handler"
	"github.com/onetwoclimb/internal/server/restapi"
	"github.com/onetwoclimb/internal/server/restapi/operations"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var config Config

var l = logrus.New()

func init() {
	RootCmd.AddCommand(cmd)
	cmd.Flags().AddFlagSet(config.Flags())
}

var cmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	RunE: func(cmd *cobra.Command, args []string) error {
		l.Info("start 12Climb server")
		defer l.Info("stop 12Climb server")

		// load embedded swagger file
		swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
		if err != nil {
			return errors.WithStack(err)
		}

		// create new service API
		api := operations.NewOneTwoClimbAPI(swaggerSpec)
		server := restapi.NewServer(api)
		defer func() {
			if err := server.Shutdown(); err != nil {
				l.WithError(err).Error()
			}
		}()

		flag.Parse()
		// set the port this service will be run on
		server.Port = config.Port
		server.Host = config.Host
		server.ReadTimeout = config.ReadTimeout
		server.WriteTimeout = config.WriteTimeout

		handler.New().ConfigureHandlers(api)

		// serve API
		if err := server.Serve(); err != nil {
			return errors.WithStack(err)
		}

		return nil
	},
}
