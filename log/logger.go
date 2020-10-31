package log

import (
	"context"
	"io"
)

var (
	standard      = newEntry(newRus()).WithField("prefix", "MAIN")
	cHook         = new(callerHook)
	textFormatter = &DefaultFormatter{callerMaxChars: 15}
)

func init() {
	SetFormatter(textFormatter)
	AddHook(cHook)
	SetLevel("trace")
	SetReportCaller(true)
}

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Printf(format string, args ...interface{})
}

func Standard() Logger {
	return standard
}

func AddHook(hook Hook) {
	standard.Logger.AddHook(hook)
}

func SetLevel(level string) {
	lv, err := ParseLevel(level)
	if err == nil {
		standard.Logger.SetLevel(lv)
	}
}

func SetOut(out io.Writer) {
	standard.Logger.SetOutput(out)
}

func SetFormatter(formatter Formatter) {
	standard.Logger.SetFormatter(formatter)
}

func SetReportCaller(reportCaller bool) {
	standard.Logger.SetReportCaller(reportCaller)
}

func SetPrefix(prefix string) {
	standard = standard.WithField("prefix", prefix)
}

func WithContext(ctx context.Context) Logger {
	return standard.WithContext(ctx)
}

func WithPrefix(prefix string) Logger {
	return standard.WithField("prefix", prefix)
}

func Debugf(format string, args ...interface{}) {
	standard.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	standard.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	standard.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	standard.Errorf(format, args...)
}

func Printf(format string, args ...interface{}) {
	standard.Logf(ErrorLevel, format, args...)
}
