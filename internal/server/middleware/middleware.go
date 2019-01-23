package middleware

import (
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/facebookgo/stack"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var l = logrus.New()

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rerr := recover(); rerr != nil {
				defer func(rw http.ResponseWriter) {
					rw.WriteHeader(http.StatusInternalServerError)
				}(w)

				httprequest, err := httputil.DumpRequest(r, false)
				if err != nil {
					l.WithError(err).Errorf("failed to dump the request on panic %v", rerr)
					return
				}

				// todo get stack trace from rerr
				l.WithError(errors.New("panic recovery")).
					WithField("stack", stack.Caller(3).String()).
					WithField("request", string(httprequest)).
					Error("api panic middleware triggered")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		// todo add request status code
		l.WithFields(logrus.Fields{
			"method":       r.Method,
			"path":         r.URL,
			"took-seconds": time.Since(start).Seconds(),
		}).Info("REST request")
	})
}
