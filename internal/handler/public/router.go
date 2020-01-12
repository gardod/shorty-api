package public

import (
	"net/http"

	"github.com/gardod/shorty-api/internal/handler/public/link"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func GetRouter() http.Handler {
	r := chi.NewRouter()

	// TODO: implement own requestID and logger middleware with logrus
	r.Use(
		middleware.DefaultCompress,
		middleware.StripSlashes,
		middleware.RequestID,
		middleware.Logger,
	)

	r.Mount("/link", link.GetRouter())

	return r
}
