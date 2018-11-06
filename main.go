package main

import (
	"os"

	"go.uber.org/zap"

	"github.com/skuid/balsa/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		zap.L().Error("Encountered an error with root cobra command", zap.Error(err))
		os.Exit(-1)
	}

}
