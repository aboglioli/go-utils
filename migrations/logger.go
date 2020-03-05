package migrations

import "fmt"

type Logger interface {
	Log(string, ...interface{})
}

type defaultLogger struct{}

func DefaultLogger() *defaultLogger {
	return &defaultLogger{}
}

func (l *defaultLogger) Log(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
