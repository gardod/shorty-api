package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gardod/shorty-api/internal/driver/http/response"
)

type Renderer func(w http.ResponseWriter, v interface{}, code int)

func Recoverer(render Renderer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil && rec != http.ErrAbortHandler {
					GetLogger(r.Context()).WithField("error", string(debug.Stack())).Error("Recovered from a panic")

					render(w, response.ErrInternal, http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
