package middleware

import (
	"net/http"
	"runtime/debug"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil && rec != http.ErrAbortHandler {
				GetLogger(r.Context()).WithField("error", string(debug.Stack())).Error("Recovered from a panic")

				// TODO: replace with own response package
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
