// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"io"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/onetwoclimb/internal/server/restapi/operations"
)

//go:generate swagger generate server --target ../internal/server --name OneTwoClimb --spec ../api/spec.yaml --exclude-main

func configureFlags(api *operations.OneTwoClimbAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.OneTwoClimbAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	api.ImagePngImageJpegProducer = runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
		return errors.NotImplemented("imagePngImageJpeg producer has not yet been implemented")
	})

	api.DelBoardColorHandler = operations.DelBoardColorHandlerFunc(func(params operations.DelBoardColorParams) middleware.Responder {
		return middleware.NotImplemented("operation .DelBoardColor has not yet been implemented")
	})
	api.DownloadFileHandler = operations.DownloadFileHandlerFunc(func(params operations.DownloadFileParams) middleware.Responder {
		return middleware.NotImplemented("operation .DownloadFile has not yet been implemented")
	})
	api.GetBoardColorsHandler = operations.GetBoardColorsHandlerFunc(func(params operations.GetBoardColorsParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetBoardColors has not yet been implemented")
	})
	api.PostBoardColorsHandler = operations.PostBoardColorsHandlerFunc(func(params operations.PostBoardColorsParams) middleware.Responder {
		return middleware.NotImplemented("operation .PostBoardColors has not yet been implemented")
	})
	api.UploadFileHandler = operations.UploadFileHandlerFunc(func(params operations.UploadFileParams) middleware.Responder {
		return middleware.NotImplemented("operation .UploadFile has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
