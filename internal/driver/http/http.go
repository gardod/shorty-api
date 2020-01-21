package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Serve(handler http.Handler) {
	logrus.Info("server starting")

	viper.SetDefault("api.port", "80")

	server := &http.Server{
		Addr:         ":" + viper.GetString("api.port"),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.WithError(err).Fatal("unable to start server")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.WithError(err).Fatal("unable to gracefully shut down server")
	}

	logrus.Info("server shut down")
}
