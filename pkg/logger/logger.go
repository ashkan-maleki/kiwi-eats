package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func Init() {
	var err error
	//Logger, err = zap.NewProduction() // or zap.NewDevelopment() for human-readable logs
	Logger, err = zap.NewDevelopment() // or zap.NewDevelopment() for human-readable logs
	if err != nil {
		panic(err)
	}
	defer func(Logger *zap.Logger) {
		err := Logger.Sync()
		if err != nil {
			panic(err)
		}
	}(Logger)
}
