package link

import (
	"net/http"

	"github.com/go-chi/chi"
)

func GetRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", list)
	r.Post("/", create)

	// TODO: adjust regex
	r.Get("/{short:[a-zA-Z0-9-]+}", getByShort)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", get)
		r.Patch("/", update)
		r.Delete("/", delete)
	})

	return r
}
