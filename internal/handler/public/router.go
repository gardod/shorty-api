package public

import (
	"net/http"

	"github.com/gardod/shorty-api/internal/driver/http/response"
	"github.com/gardod/shorty-api/internal/handler/public/link"
	mw "github.com/gardod/shorty-api/internal/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
		mw.Recoverer(response.JSON),
	)

	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		response.JSON(w, response.ErrNotFound, http.StatusNotFound)
	}))

	r.MethodNotAllowed(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		response.JSON(w, response.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	}))

	r.Mount("/link", link.GetRouter())

	return r
}
