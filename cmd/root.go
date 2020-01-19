package cmd

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "shorty-api",
	Short: "URL shortener API",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Usage(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLog)

	rootCmd.PersistentFlags().String("config", "", "Config file path")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "Set the logging level")
	rootCmd.PersistentFlags().BoolP("debug", "D", false, "Enable debug mode")
	viper.BindPFlags(rootCmd.PersistentFlags())
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
		logrus.WithError(err).Fatal("unable to read config")
	}

	if viper.GetBool("debug") {
		viper.Set("log-level", "trace")
	}
}

func initLog() {
	if viper.GetBool("debug") {
		logrus.SetReportCaller(true)
	}

	level, err := logrus.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		logrus.WithError(err).Fatal("invalid log level")
	}
	logrus.SetLevel(level)
}
