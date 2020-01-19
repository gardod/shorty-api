package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

const loggerCtxKey contextKey = "logger"

func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var logger logrus.FieldLogger = logrus.StandardLogger()

		reqID := middleware.GetReqID(r.Context())
		if reqID != "" {
			logger = logger.WithField("reqID", reqID)
		}

		r = r.WithContext(context.WithValue(r.Context(), loggerCtxKey, logger))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func GetLogger(ctx context.Context) logrus.FieldLogger {
	if logger, ok := ctx.Value(loggerCtxKey).(logrus.FieldLogger); ok {
		return logger
	}
	return logrus.StandardLogger()
}

func RequestLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		defer func() {
			GetLogger(r.Context()).WithFields(logrus.Fields{
				"method":     r.Method,
				"uri":        r.RequestURI,
				"protocol":   r.Proto,
				"remoteAddr": r.RemoteAddr,
				"status":     ww.Status(),
				"size":       ww.BytesWritten(),
				"duration":   time.Since(start),
			}).Info("Finished HTTP response")
		}()

		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}
