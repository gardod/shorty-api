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

	rootCmd.PersistentFlags().String("config", "", "config file path")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initConfig() {
	if configFile := viper.GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/shorty/")
		viper.AddConfigPath("$HOME/.shorty/")
		viper.AddConfigPath("./config/")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Fatal("unable to read config")
	}
}

func initLog() {
	if viper.GetBool("verbose") {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
