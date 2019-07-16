package public

import (
	"net/http"

	"github.com/go-chi/chi"
)

func GetRouter() http.Handler {
	r := chi.NewRouter()

	return r
}
