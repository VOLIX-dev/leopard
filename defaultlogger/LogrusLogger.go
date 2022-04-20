package defaultlogger

import "github.com/sirupsen/logrus"

type LogrusLogger struct {
}

func New() *LogrusLogger {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetLevel(logrus.DebugLevel)

	return &LogrusLogger{}
}

func (l LogrusLogger) Info(arg ...interface{}) {
	logrus.Info(arg...)
}

func (l LogrusLogger) Warning(arg ...interface{}) {
	logrus.Warning(arg...)
}

func (l LogrusLogger) Error(arg ...interface{}) {
	logrus.Error(arg...)
}

func (l LogrusLogger) Debug(arg ...interface{}) {
	logrus.Debug(arg...)
}
