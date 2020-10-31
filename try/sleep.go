package try

import (
	"context"
	"github.com/shupkg/trier/log"
	"time"
)

//幂指等待
func Sleep(ctx context.Context, sleep time.Duration, times ...int) error {
	m := 0
	if len(times) > 0 {
		m = times[0] - 1
		if m < 0 {
			m = 0
		}
	}
	sleep = time.Duration(1<<m) * sleep
	log.Debugf("sleep: %s\n", sleep.String())
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(sleep):
		return nil
	}
}

//休眠到位置
func SleepUntil(ctx context.Context, when time.Time) error {
	if when.Before(time.Now()) {
		return nil
	}
	if sleep := -time.Since(when); sleep > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleep):
			return nil
		}
	}
	return nil
}
