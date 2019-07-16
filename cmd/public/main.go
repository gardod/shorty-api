package public

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gardod/shorty-api/internal/handler/public"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Serve() {
	logrus.Debug("server starting")

	viper.SetDefault("api.public.port", "80")

	server := &http.Server{
		Addr:    ":" + viper.GetString("api.public.port"),
		Handler: public.GetRouter(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.WithError(err).Error("unable to start server")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("unable to gracefully shut down server")
	}

	logrus.Debug("server shut down")
}
