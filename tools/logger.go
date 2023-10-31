package tools

import (
	"go.uber.org/zap"
	"sync"
)

var logger zap.Logger
var once = new(sync.Once)

func Logger() *zap.Logger {
	once.Do(func() {
		newLogger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		logger = *newLogger
	})
	return &logger
}
