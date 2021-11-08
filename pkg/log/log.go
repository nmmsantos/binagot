package log

import (
	"sync"

	"go.uber.org/zap"
)

var once sync.Once
var sugar *zap.SugaredLogger

func GetLogger() *zap.SugaredLogger {
	once.Do(func() {
		if logger, err := zap.NewDevelopment(); err != nil {
			panic(err)
		} else {
			sugar = logger.Sugar()
		}
	})

	return sugar
}
