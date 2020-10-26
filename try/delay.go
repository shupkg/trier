package try

import (
	"context"
	"github.com/shupkg/trier/log"
	"time"
)

//延迟返回错误
func DelayEx(ctx context.Context, err error, base time.Duration, times ...int) error {
	if err := Delay(ctx, base, times...); err != nil {
		return err
	}
	return err
}

//幂指等待
func Delay(ctx context.Context, base time.Duration, times ...int) error {
	m := 0
	if len(times) > 0 {
		m = times[0] - 1
		if m < 0 {
			m = 0
		}
	}
	sleep := time.Duration(1<<m) * base
	log.Debugf("sleep: %s\n", sleep.String())
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(sleep):
		return nil
	}
}
