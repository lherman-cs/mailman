//+build prod

package main

import (
	"os"
	"path"

	"go.uber.org/zap"
)

var production = true
var logger *zap.SugaredLogger
var userDir string

func init() {
	// setup logger
	var cfg zap.Config

	cfg = zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"stderr",
	}
	cfg.ErrorOutputPaths = []string{
		"stderr",
	}
	// TODO! Allow users to change this level dynamically
	cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	rawLogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	logger = rawLogger.Sugar()

	// setup userDir
	userDir, err = os.UserConfigDir()
	if err != nil {
		logger.Fatal(err)
	}
	userDir = path.Join(userDir, name)
}
