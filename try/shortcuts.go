package try

import (
	"context"
	"time"
)

func RunWait(ctx context.Context, worker WorkFunc, onError ErrorFunc) error {
	return (&Trier{}).SetDelayOnError(time.Second, true).HandleError(onError).SetWork(worker).Run(ctx)
}

func Run(ctx context.Context, worker WorkFunc, onError ErrorFunc) error {
	return (&Trier{}).HandleError(onError).SetWork(worker).Run(ctx)
}
