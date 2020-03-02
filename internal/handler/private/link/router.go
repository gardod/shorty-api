package link

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gardod/shorty-api/internal/driver/http/response"
	"github.com/gardod/shorty-api/internal/model"
	"github.com/gardod/shorty-api/internal/repository"
	"github.com/gardod/shorty-api/internal/service"

	"github.com/go-chi/chi"
)

const linkCtxKey string = "shorty-api link"

func GetRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/", list)
	r.Post("/", insert)

	r.Route("/{id:[0-9]+}", func(r chi.Router) {
		r.Use(LinkCtx)

		r.Get("/", get)
		r.Patch("/", update)
		r.Delete("/", delete)
	})

	return r
}

func LinkCtx(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

		link, err := service.NewLink(ctx).GetByID(ctx, id)
		switch err {
		case nil:
		case repository.ErrNoResults:
			response.JSON(w, response.ErrNotFound, http.StatusNotFound)
			return
		default:
			response.JSON(w, response.ErrInternal, http.StatusInternalServerError)
			return
		}

		r = r.WithContext(context.WithValue(ctx, linkCtxKey, link))
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func GetLink(ctx context.Context) *model.Link {
	return ctx.Value(linkCtxKey).(*model.Link)
}
