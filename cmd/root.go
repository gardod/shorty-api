package cmd

import (
	"strings"

	"github.com/gardod/shorty-api/internal/driver/postgres"
	"github.com/gardod/shorty-api/internal/driver/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "shorty-api",
	Short: "URL shortener API",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Usage(); err != nil {
			logrus.WithError(err).Fatal("Unable to output usage instructions")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Unable to execute command")
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLog, initDrivers)

	rootCmd.PersistentFlags().String("config", "", "Config file path")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Set the logging level")
	rootCmd.PersistentFlags().BoolP("debug", "D", false, "Enable debug mode")
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		logrus.WithError(err).Fatal("Unable to bind flags")
	}
}

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if configFile := viper.GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/shorty/")
		viper.AddConfigPath("$HOME/.shorty/")
		viper.AddConfigPath("./config/")
	}
	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Fatal("Unable to read config")
	}
}

func initLog() {
	if viper.GetBool("debug") {
		logrus.SetReportCaller(true)
		viper.SetDefault("log-level", "debug")
	}

	level, err := logrus.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		logrus.WithError(err).Fatal("Invalid log level")
	}
	logrus.SetLevel(level)
}

func initDrivers() {
	postgres.Init()
	redis.Init()
}
