package es

import (
	"go.uber.org/zap"
)

type loggerES struct {
	*zap.Logger
}

func newLoggerES(log *zap.Logger) *loggerES {
	return &loggerES{
		log,
	}
}

func (l *loggerES) Printf(format string, v ...interface{}) {
	l.Sugar().Debugf(format, v...)
}
