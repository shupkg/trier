package log

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var locShanghai = time.FixedZone("Asia/Shanghai", 8*3600)

type DefaultFormatter struct {
	callerMaxChars int
	prefixMaxChars int
}

func (f *DefaultFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	var (
		entryLevel = f.formatLevel(entry)
		entryTime  = f.formatTime(entry)
		prefix     = f.formatPrefix(entry)
		caller     = f.formatCaller(entry)
	)

	lines := strings.Split(strings.Trim(entry.Message, "\n"), "\n")
	for _, line := range lines {
		b.WriteString("⇨ ")
		b.WriteString(entryLevel)
		b.WriteString(entryTime)
		b.WriteString(prefix)
		b.WriteString(caller)
		b.WriteString(line)
		b.WriteByte('\n')
	}

	if len(entry.Data) > 0 {
		b.WriteString("⇨ ")
		b.WriteString(entryLevel)
		b.WriteString(entryTime)
		b.WriteString(prefix)
		b.WriteString(caller)
		for k, v := range entry.Data {
			_, _ = fmt.Fprintf(b, " %s=", k)
			f.appendValue(b, v)
		}
		b.WriteByte('\n')
	}

	return b.Bytes(), nil
}

func (f *DefaultFormatter) formatLevel(entry *Entry) string {
	return strings.ToUpper(entry.Level.String()[:4]) + " | "
}

func (f *DefaultFormatter) formatTime(entry *Entry) string {
	return entry.Time.In(locShanghai).Format("2006-01-02/150405.000+8") + " | "
}

func (f *DefaultFormatter) formatPrefix(entry *Entry) string {
	if prefixI, ok := entry.Data["prefix"]; ok {
		if prefix, _ := prefixI.(string); prefix != "" {
			delete(entry.Data, "prefix")
			if f.prefixMaxChars < 8 {
				f.prefixMaxChars = 8
			}
			if len(prefix) > f.prefixMaxChars {
				return " .." + prefix[:f.prefixMaxChars-2] + " | "
			}
			return " " + prefix + strings.Repeat(" ", f.prefixMaxChars-len(prefix)) + " | "
		}
	}
	return " " + strings.Repeat("-", f.prefixMaxChars) + " | "
}

func (f *DefaultFormatter) formatCaller(entry *Entry) string {
	if entry.Logger != nil && entry.Caller != nil {
		if entry.Caller != nil {
			file := filepath.Base(entry.Caller.File)
			if f.callerMaxChars < 10 {
				f.callerMaxChars = 10
			}
			if l := len(file); l > f.callerMaxChars {
				file = ".." + file[l-f.callerMaxChars-2:]
			}
			file = fmt.Sprintf(" %"+strconv.Itoa(f.callerMaxChars)+"s:%-3d | ", file, entry.Caller.Line)
			return file
		}
	}
	return ""
}

func (f *DefaultFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

func (f *DefaultFormatter) needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}
