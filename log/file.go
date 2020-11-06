package log

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

func File(o FileOption) io.Writer {
	if o.Filename != "" {
		if fn, _ := filepath.Abs(o.Filename); fn != "" {
			_ = os.MkdirAll(filepath.Dir(fn), 0755)
		}
	}
	return &lumberjack.Logger{
		Filename:   o.Filename,
		MaxSize:    o.MaxSize,
		MaxAge:     o.MaxAge,
		MaxBackups: o.MaxBackups,
		LocalTime:  o.LocalTime,
		Compress:   o.Compress,
	}
}

func Multi(o FileOption) io.Writer {
	return io.MultiWriter(Terminal(), File(o))
}

func Terminal() io.Writer {
	return os.Stdout
}
