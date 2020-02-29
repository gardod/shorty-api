package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Serve(handler http.Handler) {
	logrus.Info("Server starting")

	viper.SetDefault("api.port", "80")

	h2s := &http2.Server{
		IdleTimeout: time.Second * 60,
	}

	h1s := &http.Server{
		Addr:         ":" + viper.GetString("api.port"),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      h2c.NewHandler(handler, h2s),
	}

	go func() {
		if err := h1s.ListenAndServe(); err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("Unable to start server")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h1s.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("Unable to gracefully shut down server")
	}

	logrus.Info("Server shut down")
}
