package link

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gardod/shorty-api/internal/driver/http/response"
	"github.com/gardod/shorty-api/internal/service"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v3"
)

func list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	from := time.Now()
	if param := q.Get("from"); param != "" {
		var err error
		from, err = time.Parse(time.RFC3339, param)
		if err != nil {
			response.JSON(w, validation.Errors{"from": err}, http.StatusBadRequest)
			return
		}
	}

	limit := 100
	if param := q.Get("limit"); param != "" {
		var err error
		limit, err = strconv.Atoi(param)
		if err != nil {
			response.JSON(w, validation.Errors{"limit": err}, http.StatusBadRequest)
			return
		}
	}

	link, err := service.NewLink(ctx).Get(ctx, from, limit)
	if err != nil {
		response.JSON(w, response.ErrInternal, http.StatusInternalServerError)
		return
	}

	response.JSON(w, link, http.StatusOK)
}

func insert(w http.ResponseWriter, r *http.Request) {}

func get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	link, err := service.NewLink(ctx).GetByID(ctx, id)
	switch err {
	case nil:
	case sql.ErrNoRows:
		response.JSON(w, response.ErrNotFound, http.StatusNotFound)
		return
	default:
		response.JSON(w, response.ErrInternal, http.StatusInternalServerError)
		return
	}

	response.JSON(w, link, http.StatusOK)
}

func update(w http.ResponseWriter, r *http.Request) {}

func delete(w http.ResponseWriter, r *http.Request) {}
