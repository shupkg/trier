package try

import (
	"context"
	"errors"
	"time"
)

var ErrUserBreak = errors.New("user break")

type TxFunc func(ctx context.Context, runTimes int) (next time.Time, err error)

func RunTx(ctx context.Context, txFunc TxFunc) error {
	for i := 0; ; i++ {
		if next, err := txFunc(ctx, i+1); err != nil {
			if errors.Is(err, ErrUserBreak) {
				return nil
			}
			return err
		} else if err := When(ctx, next); err != nil {
			return err
		}
	}
}

func When(ctx context.Context, when time.Time) error {
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
