package log

import (
	"github.com/sirupsen/logrus"
)

type (
	Rus   = logrus.Logger
	Entry = logrus.Entry
	Level = logrus.Level
)

var (
	newEntry = logrus.NewEntry
	newRus   = logrus.New

	ParseLevel = logrus.ParseLevel
	AllLevels  = logrus.AllLevels
)

const (
	TraceLevel = logrus.TraceLevel
	DebugLevel = logrus.DebugLevel
	InfoLevel  = logrus.InfoLevel
	WarnLevel  = logrus.WarnLevel
	ErrorLevel = logrus.ErrorLevel
	FatalLevel = logrus.FatalLevel
	PanicLevel = logrus.PanicLevel
)

type Hook interface {
	Levels() []Level
	Fire(*Entry) error
}

type Formatter interface {
	Format(*Entry) ([]byte, error)
}
