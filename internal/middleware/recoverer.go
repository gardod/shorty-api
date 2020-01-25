package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gardod/shorty-api/internal/driver/http/response"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil && rec != http.ErrAbortHandler {
				GetLogger(r.Context()).WithField("error", string(debug.Stack())).Error("Recovered from a panic")

				response.JSON(w, response.ErrInternal, http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
