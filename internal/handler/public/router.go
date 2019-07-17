package public

import (
	"net/http"

	"github.com/gardod/shorty-api/internal/handler/public/link"
	"github.com/go-chi/chi"
)

func GetRouter() http.Handler {
	r := chi.NewRouter()

	r.Mount("link", link.GetRouter())

	return r
}
