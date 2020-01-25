package link

import (
	"net/http"

	"github.com/go-chi/chi"
)

func GetRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/{short:[a-zA-Z0-9-]+}", getByShort)

	return r
}
