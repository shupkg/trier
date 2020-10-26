package log

import (
	"runtime"
	"strings"
	"sync"
)

const (
	rusPackage = "github.com/sirupsen/logrus"
)

var frameNil = runtime.Frame{}

type callerHook struct {
	once       sync.Once
	logPackage string
	stackStart int
}

func (c *callerHook) Levels() []Level {
	return AllLevels
}

func (c *callerHook) Fire(entry *Entry) error {
	entry.Caller = c.getCaller()
	return nil
}

func (c *callerHook) getCaller() *runtime.Frame {
	c.once.Do(func() {
		pcs := make([]uintptr, 25)
		depth := runtime.Callers(0, pcs)
		for i := 0; i < depth; i++ {
			if funcName := runtime.FuncForPC(pcs[i]).Name(); strings.Contains(funcName, "getCaller") {
				c.logPackage = c.getPackageName(funcName)
				c.stackStart = i + 1
				break
			}
		}
	})

	pcs := make([]uintptr, 25)
	frames := runtime.CallersFrames(pcs[:runtime.Callers(c.stackStart, pcs)])
	var caller runtime.Frame
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if strings.HasPrefix(frame.Function, c.logPackage) || strings.HasPrefix(frame.Function, rusPackage) {
			caller = frameNil
			continue
		}
		if caller == frameNil {
			caller = frame
		}
	}

	if caller.Func != nil {
		return &caller
	}
	return nil
}

func (c *callerHook) getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

var _ Hook = (*callerHook)(nil)
