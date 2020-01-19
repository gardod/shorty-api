package public

import (
	"net/http"

	"github.com/gardod/shorty-api/internal/handler/public/link"
	mw "github.com/gardod/shorty-api/internal/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func GetRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.DefaultCompress,
		middleware.StripSlashes,
		middleware.RealIP,
		middleware.RequestID,
		mw.Logger,
		mw.RequestLogger,
		mw.Recoverer,
	)

	r.Mount("/link", link.GetRouter())

	return r
}
