package middleware

import (
	"context"
	"net/http"

	"github.com/gardod/shorty-api/internal/driver/postgres"
)

const dbCtxKey contextKey = "db"

func Database(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		client := postgres.NewClient()

		r = r.WithContext(context.WithValue(r.Context(), dbCtxKey, client))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func GetDB(ctx context.Context) *postgres.Client {
	return ctx.Value(dbCtxKey).(*postgres.Client)
}
