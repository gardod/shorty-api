package link

import (
	"net/http"

	"github.com/go-chi/chi"
)

func GetRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", list)
	r.Post("/", insert)

	r.Route("/{id:[0-9]+}", func(r chi.Router) {
		r.Get("/", get)
		r.Put("/", update)
		r.Delete("/", delete)
	})

	return r
}
