package link

import (
	"net/http"

	"github.com/gardod/shorty-api/internal/driver/http/response"
	"github.com/gardod/shorty-api/internal/repository"
	"github.com/gardod/shorty-api/internal/service"

	"github.com/go-chi/chi"
)

func getByShort(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	short := chi.URLParam(r, "short")

	link, err := service.NewLink(ctx).GetByShort(ctx, short)
	switch err {
	case nil:
	case repository.ErrNoResults:
		response.JSON(w, response.ErrNotFound, http.StatusNotFound)
		return
	default:
		response.JSON(w, response.ErrInternal, http.StatusInternalServerError)
		return
	}

	response.JSON(w, link, http.StatusOK)
}
