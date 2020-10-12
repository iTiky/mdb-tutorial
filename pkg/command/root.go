package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/itiky/mdb-tutorial/pkg/common"
)

var (
	configFile string
	logLevel   string
)

// rootCmd is a base command.
var rootCmd = &cobra.Command{
	Use:   "commands",
	Short: "MongoDB tutorial",
}

// Execute starts rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads application config (path is taken from viper).
func initConfig() {
	if configFile == "" {
		logrus.Infof("Config file: not used")
		return
	}

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Config file (%s): read failed: %v", configFile, err)
	}
	logrus.Infof("Config file: %s", viper.ConfigFileUsed())
}

// initLogger creates a new logger (logging level is taken from viper).
func initLogger() *logrus.Logger {
	l := logrus.New()

	logLvlStr := viper.GetString(common.AppLogLevel)
	logLvl, err := logrus.ParseLevel(logLvlStr)
	if err != nil {
		logrus.Fatalf("Logger init: invalid logging level: %s", logLvlStr)
	}

	l.SetLevel(logLvl)
	logrus.Warnf("Logging level: %s", logLvl)

	return l
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file path")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "logging level (debug, info, warn, error, fatal, panic)")
	if err := viper.BindPFlag(common.AppLogLevel, rootCmd.PersistentFlags().Lookup("log-level")); err != nil {
		panic(err)
	}

	cobra.OnInitialize(initConfig)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("MDB_TUTORIAL")
	viper.AutomaticEnv()
}
