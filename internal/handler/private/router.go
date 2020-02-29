package private

import (
	"net/http"

	"github.com/gardod/shorty-api/internal/driver/http/response"
	"github.com/gardod/shorty-api/internal/handler/private/link"
	mw "github.com/gardod/shorty-api/internal/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
)

func GetRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.StripSlashes,
		middleware.RealIP,
		middleware.RequestID,
		mw.Logger,
		mw.RequestLogger,
		mw.Database,
		mw.Cache,
		mw.Recoverer(response.Gob),
	)

	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		response.Gob(w, response.ErrNotFound, http.StatusNotFound)
	}))

	r.MethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		response.Gob(w, response.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}))

	if viper.GetBool("debug") {
		r.Mount("/debug", mw.Profiler())
	}

	r.Mount("/link", link.GetRouter())

	return r
}
