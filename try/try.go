package try

import (
	"context"
	"errors"
	"time"

	"github.com/shupkg/trier/log"
)

var ErrStopTry = errors.New("stop try")

var sleeps = []int{1, 1, 2, 4, 8, 16, 32, 32, 32}

type WorkFunc = func(tryTime int) (err error)

type ErrorFunc = func(tryTime int, ex error) (err error)

func Try(ctx context.Context, tryFunc WorkFunc, waits ...Wait) error {
	var wait Wait
	if len(waits) > 0 {
		wait = waits[0]
	} else {
		wait = SleepWait(time.Second)
	}

	for i := 1; ; i++ {
		if done, err := isDone(tryFunc(i)); done {
			return err
		}

		if wait != nil {
			//额外增加休眠
			var (
				sleepWhen = time.Now()
				b         = wait.Burst()
				n         = 0
			)

			if i < len(sleeps) {
				n = sleeps[i]
			} else {
				n = sleeps[len(sleeps)-1]
			}

			if n <= b {
				if err := wait.WaitN(ctx, n); err != nil {
					return err
				}
			} else {
				for j := n; j > 0; j -= b {
					k := j
					if k > b {
						k = b
					}
					if err := wait.WaitN(ctx, k); err != nil {
						return err
					}
				}
			}
			log.Debugf("休眠: %s", time.Since(sleepWhen).String())
		}
	}
}

func isDone(err error) (bool, error) {
	//正常执行或者用户取消重试
	if err == nil || errors.Is(err, ErrStopTry) {
		return true, nil
	}
	//上下文终止，退出重试
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true, err
	}
	return false, err
}
