package try

import (
	"context"
	"math"
	"time"
)

type Wait interface {
	Wait(ctx context.Context) error
	WaitN(ctx context.Context, n int) error
	Burst() int
}

func SleepWait(base time.Duration) Wait {
	return sleepWait(base)
}

type sleepWait time.Duration

func (l sleepWait) Wait(ctx context.Context) error {
	return Sleep(ctx, time.Second)
}

func (l sleepWait) WaitN(ctx context.Context, n int) error {
	return Sleep(ctx, time.Second*time.Duration(n))
}

func (l sleepWait) Burst() int {
	return math.MaxInt32
}

var _ Wait = sleepWait(0)
