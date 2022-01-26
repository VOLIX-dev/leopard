package leopard

import "leopard/defaultlogger"

var Logger LoggerInterface

func init() {
	Logger = defaultlogger.New()
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Warning(args ...interface{}) {
	Logger.Warning(args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

type LoggerInterface interface {
	Info(arg ...interface{})
	Warning(arg ...interface{})
	Error(arg ...interface{})
	Debug(arg ...interface{})
}
