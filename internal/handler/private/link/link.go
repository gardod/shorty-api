package link

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gardod/shorty-api/internal/driver/http/response"
	"github.com/gardod/shorty-api/internal/model"
	"github.com/gardod/shorty-api/internal/repository"
	"github.com/gardod/shorty-api/internal/service"

	vld "github.com/go-ozzo/ozzo-validation/v4"
)

func list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	from := time.Now()
	if param := q.Get("from"); param != "" {
		var err error
		from, err = time.Parse(time.RFC3339, param)
		if err != nil {
			response.JSON(w, vld.Errors{"from": err}, http.StatusBadRequest)
			return
		}
	}

	limit := 100
	if param := q.Get("limit"); param != "" {
		var err error
		limit, err = strconv.Atoi(param)
		if err != nil {
			response.JSON(w, vld.Errors{"limit": err}, http.StatusBadRequest)
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

func insert(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	link := &model.Link{}

	if err := json.NewDecoder(r.Body).Decode(&link.LinkRequest); err != nil {
		response.JSON(w, response.ErrParse.WithDetails(err), http.StatusBadRequest)
		return
	}

	err := service.NewLink(ctx).Insert(ctx, link)
	if err != nil {
		if _, ok := err.(vld.Errors); ok {
			response.JSON(w, response.ErrValidation.WithDetails(err), http.StatusUnprocessableEntity)
			return
		}
		response.JSON(w, response.ErrInternal, http.StatusInternalServerError)
		return
	}

	response.JSON(w, link, http.StatusOK)
}

func get(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, GetLink(r.Context()), http.StatusOK)
}

func update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	link := GetLink(ctx)

	if err := json.NewDecoder(r.Body).Decode(&link.LinkRequest); err != nil {
		response.JSON(w, response.ErrParse.WithDetails(err), http.StatusBadRequest)
		return
	}

	err := service.NewLink(ctx).Update(ctx, link)
	switch err {
	case nil:
	case repository.ErrNoResults:
		response.JSON(w, response.ErrNotFound, http.StatusNotFound)
		return
	default:
		if _, ok := err.(vld.Errors); ok {
			response.JSON(w, response.ErrValidation.WithDetails(err), http.StatusUnprocessableEntity)
			return
		}
		response.JSON(w, response.ErrInternal, http.StatusInternalServerError)
		return
	}

	response.JSON(w, link, http.StatusOK)
}

func delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	link := GetLink(ctx)

	err := service.NewLink(ctx).Delete(ctx, link)
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
