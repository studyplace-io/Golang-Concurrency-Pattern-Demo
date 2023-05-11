package cronfunc

import (
	"context"
	"k8s.io/apimachinery/pkg/util/runtime"
	"math/rand"
	"time"
)

// NeverStop 定义的一个chan
var NeverStop <-chan struct{} = make(chan struct{})

// NeverStopFunc 定时执行f()，不会退出
func NeverStopFunc(f func(), period time.Duration) {
	CronUntil(f, period, NeverStop)
}

// CronUntil 定时轮巡执行f()，定时周期是在第一次调用后才开始计算
func CronUntil(f func(), period time.Duration, stopCh <-chan struct{}) {
	LoopWithRandomFactor(f, period, 0.0, true, stopCh)
}

// CronUntilWithContext 定时轮巡执行f()，定时周期是在第一次调用后才开始计算，加入ctx
func CronUntilWithContext(ctx context.Context, f func(context.Context), period time.Duration) {
	CronUntilWithContextWithRandomFactor(ctx, f, period, 0.0, true)
}

// CronUntilWithContextWithRandomFactor 定时轮巡执行f(ctx)，定时周期是在第一次调用后才开始计算，
// 并加入随机因子让轮巡执行的时间不固定
func CronUntilWithContextWithRandomFactor(ctx context.Context, f func(context.Context), period time.Duration, randomFactor float64, sliding bool) {
	LoopWithRandomFactor(func() { f(ctx) }, period, randomFactor, sliding, ctx.Done())
}

// CronUntilWithContextBefore 定时轮巡执行f()，定时周期是在第一次调用前开始计算，加入ctx
func CronUntilWithContextBefore(ctx context.Context, f func(context.Context), period time.Duration) {
	CronUntilWithContextWithRandomFactor(ctx, f, period, 0.0, false)
}

// CronUntilBefore 定时轮巡执行f()，定时周期是在第一次调用前开始计算
func CronUntilBefore(f func(), period time.Duration, stopCh <-chan struct{}) {
	LoopWithRandomFactor(f, period, 0.0, false, stopCh)
}

// LoopWithRandomFactor 定时执行f()的主要逻辑，for循环不断轮巡执行，并使用定时器来实现定时调用
func LoopWithRandomFactor(f func(), period time.Duration, randomFactor float64, sliding bool, stopCh <-chan struct{}) {
	var t *time.Timer
	var sawTimeout bool

	for {
		// 最前面需要先检查是否退出
		select {
		case <-stopCh:
			return
		default:
		}

		// 随机因子
		randomPeriod := period
		if randomFactor > 0.0 {
			randomPeriod = addRandomFactor(period, randomFactor)
		}

		// 不是sliding，就是先计时，再执行f()
		if !sliding {
			t = resetOrReuseTimer(t, randomPeriod, sawTimeout)
		}

		// 执行 f()
		func() {
			defer runtime.HandleCrash()
			f()
		}()

		// 如果是sliding，就是先执行，再开始计时
		if sliding {
			t = resetOrReuseTimer(t, randomPeriod, sawTimeout)
		}

		// 执行到这里会阻塞，如果定时器到了就变更sawTimeout为true，
		// 如果收到通知信号，退出
		select {
		case <-stopCh:
			return
		case <-t.C:
			sawTimeout = true
		}
	}
}

// 加入随机因子
func addRandomFactor(duration time.Duration, maxFactor float64) time.Duration {
	if maxFactor <= 0.0 {
		maxFactor = 1.0
	}
	wait := duration + time.Duration(rand.Float64()*maxFactor*float64(duration))
	return wait
}

// 重新初始化计时器
func resetOrReuseTimer(t *time.Timer, d time.Duration, sawTimeout bool) *time.Timer {
	if t == nil {
		return time.NewTimer(d)
	}
	// 如果没有停止且timeout了，
	if !t.Stop() && !sawTimeout {
		<-t.C
	}
	t.Reset(d)
	return t
}
