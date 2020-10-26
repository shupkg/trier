package try

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestDelay(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)
	err := Run(ctx, func(tryTime int) error {
		fmt.Printf("第%d次尝试\n", tryTime)
		return DelayEx(ctx, errors.New("自己的错误"), time.Second*3)
	}, func(_ *Trier, tryTime int, err error) error {
		t.Log(err)
		return err
	})
	t.Log(err)
	//for i := 0; i < 20; i++ {
	//	if err := Delay(ctx, time.Second, i); err != nil {
	//		t.Log("do...", err)
	//	}
	//}
}
