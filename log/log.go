package log

import (
	l "github.com/sirupsen/logrus"
)

var level int = 0

const (
	INFO    int = 0
	WARNING int = 1
	ERROR   int = 2
)

func SetLevel(le int) {
	level = le
	switch level {
	case INFO:
		l.SetLevel(l.InfoLevel)
	case WARNING:
		l.SetLevel(l.WarnLevel)
	case ERROR:
		l.SetLevel(l.ErrorLevel)
	}
}

func Info(params ...interface{}) {
	if level > INFO {
		return
	}
	l.Info(params...)
}

func Warning(params ...interface{}) {
	if level > WARNING {
		return
	}
	l.Warning(params...)
}

func Error(params ...interface{}) {
	if level > ERROR {
		return
	}
	l.Error(params...)
}
