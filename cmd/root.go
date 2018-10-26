package cmd

import (
	"os"

	"github.com/skuid/condparse/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var sugar = zap.L().Sugar()

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "CondParse",
	Short:   "A utility and set of functions to parse condition logic",
	Version: version.Name,
}

// Execute is used as entrypoint to the cobra commands
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		sugar.Error("encountered an error on root command execution", zap.Error(err))
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().Bool("debug", false, "Debug Mode Switch")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if err := viper.BindPFlags(RootCmd.PersistentFlags()); err != nil {
		sugar.Error("encountered an error on viper flag binding", zap.Error(err))
		os.Exit(1)
	}

	viper.SetEnvPrefix("condparse")
	viper.AutomaticEnv()
}
