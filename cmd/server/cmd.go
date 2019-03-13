package server

import (
	"flag"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/onetwoclimb/cmd/config"
	"github.com/onetwoclimb/internal/server/handler"
	serverMW "github.com/onetwoclimb/internal/server/middleware"
	"github.com/onetwoclimb/internal/server/restapi"
	"github.com/onetwoclimb/internal/server/restapi/operations"
	"github.com/onetwoclimb/internal/storages"
)

var cfg config.Config

var l = logrus.New()

func init() {
	Cmd.Flags().AddFlagSet(cfg.Flags())
}

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "run server",
	RunE: func(cmd *cobra.Command, args []string) error {
		flag.Parse()

		if err := cfg.Validate(); err != nil {
			return errors.Wrap(err, "storage folder doesn't exists")
		}

		l.Info("start 12Climb server")
		defer l.Info("stop 12Climb server")

		mysql, err := storages.New(&cfg.DB)
		if err != nil {
			return errors.Wrap(err, "failed to init mySQL db")
		}
		storage, err := storages.NewMySQLStorage(mysql)
		if err != nil {
			return errors.Wrap(err, "failed to init storage")
		}

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

		// set the port this service will be run on
		server.Port = cfg.Port
		server.Host = cfg.Host
		server.ReadTimeout = cfg.ReadTimeout
		server.WriteTimeout = cfg.WriteTimeout

		handler.New(storage, cfg).ConfigureHandlers(api)

		server.SetHandler(serverMW.PanicRecovery(serverMW.Logger(api.Serve(middleware.PassthroughBuilder))))

		// serve API
		if err := server.Serve(); err != nil {
			return errors.WithStack(err)
		}

		return nil
	},
}
