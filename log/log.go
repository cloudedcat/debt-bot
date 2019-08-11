package log

import "go.uber.org/zap"

type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	IfErrorw(err error, msg string, keysAndValues ...interface{})
}

// NewZapLogger returns wrapper for zap logger
func NewZapLogger() Logger {
	logger, _ := zap.NewProduction()
	return zapWrapper{logger.Sugar()}
}

type zapWrapper struct{ *zap.SugaredLogger }

func (wrp zapWrapper) IfErrorw(err error, msg string, keysAndValues ...interface{}) {
	if err == nil {
		return
	}
	keysAndValues = append(keysAndValues, "error", err.Error())
	wrp.Errorw(msg, keysAndValues...)
}
