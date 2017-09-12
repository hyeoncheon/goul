package goul

import (
	"github.com/sirupsen/logrus"
)

// constants
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "wan"
	LogLevelErr   = "error"
)

// Logger is an interface of the logging facility
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}

// NewLogger returns new logger. currently an instance of logrus.
func NewLogger(level string) Logger {
	logger := logrus.New()
	lev, err := logrus.ParseLevel(level)
	if err == nil {
		logger.SetLevel(lev)
	}
	return logger
}
