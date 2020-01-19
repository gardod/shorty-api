package middleware

import (
	"context"
	"net/http"

	"github.com/gardod/shorty-api/internal/driver/redis"
)

const cacheCtxKey contextKey = "cache"

func Cache(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		client := redis.NewClient(r.Context())

		r = r.WithContext(context.WithValue(r.Context(), cacheCtxKey, client))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func GetCache(ctx context.Context) *redis.Client {
	return ctx.Value(cacheCtxKey).(*redis.Client)
}
