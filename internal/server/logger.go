package server

import (
	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)

	var err error
	config := zap.NewProductionConfig()
	config.Development = true
	config.Level = atom

	//required if further abstracted
	//zapLogger, err := config.Build(zap.AddCallerSkip(1))
	zapLogger, err := config.Build()

	if err != nil {
		panic(err)
	}
	return zapLogger
}
