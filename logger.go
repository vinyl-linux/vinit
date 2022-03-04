package main

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
)

func init() {
	var err error

	if logger == nil {
		logger, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	}

	sugar = logger.Sugar()
}
