package log

import "go.uber.org/zap"

type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
}

// NewZapLogger returns wrapper for zap logger
func NewZapLogger() Logger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}
