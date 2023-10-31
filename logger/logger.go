package logger

type Logger interface {
	Error(i ...interface{})
	Errorf(format string, args ...interface{})
	Warn(i ...interface{})
	Warnf(format string, args ...interface{})
	Info(i ...interface{})
	Infof(format string, args ...interface{})
}
