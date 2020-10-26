package try

import (
	"context"
	"errors"
	"time"

	"golang.org/x/time/rate"
)

var ErrStopTry = errors.New("stop try")

type (
	ErrorFunc = func(runner *Trier, tryTime int, err error) error
	WorkFunc  = func(tryTime int) error
)

func New() *Trier {
	return &Trier{}
}

//重试器
type Trier struct {
	sleepBase       time.Duration
	sleepIdempotent bool
	errorFunc       func(runner *Trier, tryTime int, err error) error
	workerFunc      func(tryTime int) error
	rate            *rate.Limiter
}

//通过限流方法来限制重试速度
func (r *Trier) Limit(rl *rate.Limiter) *Trier {
	r.rate = rl
	return r
}

//设置错误处理函数
func (r *Trier) HandleError(onError ErrorFunc) *Trier {
	r.errorFunc = onError
	return r
}

//设置出错时延迟等待参数 base 基础等待时间， idempotent 根据重试次数幂指等待
func (r *Trier) SetDelayOnError(base time.Duration, idempotent bool) *Trier {
	r.sleepBase = base
	r.sleepIdempotent = idempotent
	return r
}

//设置执行方法
func (r *Trier) SetWork(worker WorkFunc) *Trier {
	r.workerFunc = worker
	return r
}

//执行
func (r *Trier) Run(ctx context.Context) error {
	for tryTime := 1; ; tryTime++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if r.rate != nil {
			if err := r.rate.Wait(ctx); err != nil {
				return err
			}
		}
		if err := r.workerFunc(tryTime); err != nil {
			if r.isCtxStop(err) {
				return err
			}

			if r.errorFunc != nil {
				err = r.errorFunc(r, tryTime, err)
			}

			if r.isCtxStop(err) {
				return err
			}

			if r.stopByUser(err) {
				return nil
			}

			if r.sleepBase > 0 {
				sleepTime := 0
				if r.sleepIdempotent {
					sleepTime = tryTime
				}
				if err := Delay(ctx, r.sleepBase, sleepTime); err != nil {
					return err
				}
			}
			continue
		}
		break
	}
	return nil
}

func (r *Trier) isCtxStop(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func (r *Trier) stopByUser(err error) bool {
	return errors.Is(err, ErrStopTry)
}
