package goul

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
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

// Error ...
func Error(logger Logger, module, fmt string, args ...interface{}) {
	if logger != nil {
		module = color.YellowString(module)
		logger.Errorf("["+module+"] "+fmt, args...)
	}
}

// Log ...
func Log(logger Logger, module, format string, args ...interface{}) {
	if logger != nil {
		header := fmt.Sprintf("[%v-%03d] ", module, GoID())
		message := fmt.Sprintf(header+format, args...)
		switch {
		case strings.Contains(message, "close"):
			message = color.RedString(message)
		case strings.Contains(message, "exit"):
			message = color.RedString(message)
		case strings.Contains(message, "READ"):
			message = color.GreenString(message)
		case strings.Contains(message, "start"):
			message = color.GreenString(message)
		case strings.Contains(message, "work"):
			message = color.BlueString(message)
		default:
			message = color.MagentaString(message)
		}
		logger.Debug(message)
	}
}
