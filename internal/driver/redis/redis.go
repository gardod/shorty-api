package redis

import (
	"crypto/tls"

	redis "github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var pool *redis.Client

func Init() {
	if !viper.IsSet("cache") {
		logrus.Info("Cache not enabled, skipping init")
		return
	}

	logrus.Info("Setting up cache")

	opts := redis.Options{
		Addr:     viper.GetString("cache.host") + ":" + viper.GetString("cache.port"),
		Password: viper.GetString("cache.password"),
		DB:       viper.GetInt("cache.db"),
	}
	if viper.GetBool("cache.tls") {
		opts.TLSConfig = &tls.Config{}
	}

	client := redis.NewClient(&opts)

	if err := client.Ping().Err(); err != nil {
		logrus.WithError(err).Fatal("Unable to ping cache")
	}

	pool = client
}
