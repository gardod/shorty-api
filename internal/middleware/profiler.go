package middleware

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Profiler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.NoCache)

	r.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/pprof/profile", pprof.Profile)
	r.HandleFunc("/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/pprof/trace", pprof.Trace)
	r.HandleFunc("/pprof/*", pprof.Index)
	r.HandleFunc("/pprof", pprof.Index)
	r.Handle("/vars", expvar.Handler())

	return r
}
