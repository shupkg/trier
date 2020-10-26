package log

import (
	"fmt"
	"io"
)

type Option struct {
	Level          string     `json:"level" yaml:"level" toml:"level"`
	Caller         bool       `json:"caller" yaml:"caller" toml:"caller"`
	Terminal       bool       `json:"terminal" yaml:"terminal" toml:"terminal"`
	CallerMaxChars int        `json:"caller_max_chars" yaml:"caller_max_chars" toml:"caller_max_chars"`
	File           FileOption `json:"file" yaml:"file" toml:"file"`
}

type FileOption struct {
	Filename   string `json:"filename" yaml:"filename" toml:"filename"`
	MaxSize    int    `json:"max_size" yaml:"max_size" toml:"max_size"`
	MaxAge     int    `json:"max_age" yaml:"max_age" toml:"max_age"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups" toml:"max_backups"`
	LocalTime  bool   `json:"local_time" yaml:"local_time" toml:"local_time"`
	Compress   bool   `json:"compress" yaml:"compress" toml:"compress"`
}

func (c *Option) Default() {
	c.Terminal = true
	c.Caller = true
	c.Level = "trace"

	c.File.Filename = ""
	c.File.MaxSize = 100
	c.File.MaxAge = 3
	c.File.MaxBackups = 40
	c.File.LocalTime = true
	c.File.Compress = false
}

func Apply(opt Option) {
	fmt.Printf("%-15s %s\n", "Level:", opt.Level)
	fmt.Printf("%-15s %t\n", "Caller:", opt.Caller)
	fmt.Printf("%-15s %t\n", "Terminal:", opt.Terminal)
	fmt.Printf("%-15s %d\n", "CallerMaxChars:", opt.CallerMaxChars)
	fmt.Printf("%-15s %s\n", "Filename:", opt.File.Filename)
	fmt.Printf("%-15s %dM\n", "MaxSize:", opt.File.MaxSize)
	fmt.Printf("%-15s %d天\n", "MaxAge:", opt.File.MaxAge)
	fmt.Printf("%-15s %d个\n", "MaxBackups:", opt.File.MaxBackups)
	fmt.Printf("%-15s %t\n", "LocalTime:", opt.File.LocalTime)
	fmt.Printf("%-15s %t\n", "Compress:", opt.File.Compress)
	if opt.CallerMaxChars > 0 {
		textFormatter.callerMaxChars = opt.CallerMaxChars
	}

	var out io.Writer
	if opt.File.Filename != "" {
		if opt.Terminal {
			out = Multi(opt.File)
		} else {
			out = File(opt.File)
		}
	} else {
		out = Terminal()
	}

	standard.Logger.SetOutput(out)
	if opt.Level != "" {
		SetLevel(opt.Level)
	}

	SetReportCaller(opt.Caller)
}
