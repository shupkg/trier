package signal

import (
	"github.com/shupkg/trier/log"
	"os"
	"os/signal"
	"sync"
)

type HandlerFunc = func()

var (
	signals    = make(chan os.Signal)
	handlerMap = map[os.Signal][]HandlerFunc{}
	locker     sync.Mutex
	watching   bool
)

func init() {
	go func() {
		for sig := range signals {
			if handlers, find := handlerMap[sig]; find {
				for _, handlerFunc := range handlers {
					go handlerFunc()
				}
			}
		}
	}()
}

func Handle(signal os.Signal, handlerFunc HandlerFunc) {
	handlerMap[signal] = append(handlerMap[signal], handlerFunc)
}

func Watch() {
	locker.Lock()
	defer locker.Unlock()
	if watching {
		return
	}
	watching = true
	var ss []os.Signal
	for s := range handlerMap {
		ss = append(ss, s)
		log.Debugf("监听信号: %s", s.String())
	}
	signal.Notify(signals, ss...)
}

func HandleWatch(handlerFunc HandlerFunc, signals ...os.Signal) {
	for _, o := range signals {
		Handle(o, handlerFunc)
	}
	Watch()
}
