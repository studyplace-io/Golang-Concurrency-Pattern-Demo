package forever_mode

import (
	"context"
	"time"
)

// RunForeverWithTime 永远执行 fn 直到 times 次数停止
func RunForeverWithTime(interval time.Duration, fn runFn, times int) {

	i := 0

	timer := time.NewTimer(interval)
	defer func() {
		timer.Stop()
	}()

	// 立即调用
	if err := fn(); err != nil {
		panic(err)
	}
	i++

	for {
		if i == times {
			return
		}
		select {
		case <-timer.C:
			if err := fn(); err != nil {
				panic(err)
			}
			timer.Reset(interval)
		}
		i++
	}
}

// RunWithTimeWithChannel 执行 fn，当外部关闭 stopC chan 会退出 直到 times 次数停止
func RunWithTimeWithChannel(interval time.Duration, fn runFn, stopC chan struct{}, times int) {

	i := 0

	timer := time.NewTimer(interval)
	defer func() {
		timer.Stop()
	}()

	if err := fn(); err != nil {
		panic(err)
	}

	i++

	for {
		if i == times {
			return
		}
		select {
		case <-timer.C:
			if err := fn(); err != nil {
				panic(err)
			}
			timer.Reset(interval)
		case <-stopC:
			return
		}
		i++
	}
}

// RunWithTimeWithContext 调用 fn，外部传入的 ctx 能够自由关闭 直到 times 次数停止
func RunWithTimeWithContext(interval time.Duration, fn runFn, ctx context.Context, times int) {
	i := 0
	timer := time.NewTimer(interval)
	defer func() {
		timer.Stop()
	}()

	if err := fn(); err != nil {
		panic(err)
	}
	i++
	for {
		if i == times {
			return
		}
		select {
		case <-timer.C:
			if err := fn(); err != nil {
				panic(err)
			}
			timer.Reset(interval)
		case <-ctx.Done():
			return
		}
		i++
	}
}
