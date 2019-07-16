package cmd

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "shorty-api",
	Short: "URL shortener API",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

var (
	configFile string
	verbose    bool
)

func init() {
	cobra.OnInitialize(initConfig, initLog)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is config.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("./config/")
		viper.AddConfigPath("/etc/shorty")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Error("unable to read config")
		os.Exit(1)
	}
}

func initLog() {
	if verbose {
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
