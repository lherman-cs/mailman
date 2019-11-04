//+build !prod

package main

import (
	"go.uber.org/zap"
)

var production = false
var logger *zap.SugaredLogger
var userDir = ""

func init() {
	var cfg zap.Config

	cfg = zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		"stderr",
	}
	cfg.ErrorOutputPaths = []string{
		"stderr",
	}
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	rawLogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	logger = rawLogger.Sugar()
}
