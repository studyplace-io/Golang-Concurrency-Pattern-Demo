package workqueue

import (
	"golang.org/x/time/rate"
	"time"
)

// RateLimiterQueue 限速队列接口
type RateLimiterQueue interface {
	// DelayQueue 继承 DelayQueue Interface 延迟队列基本功能
	DelayQueue
	// AddRateLimited 加入限速队列
	AddRateLimited(item interface{})
	Close()
}

// AddRateLimited 加入限速队列
func (q *rateLimitingQueue) AddRateLimited(item interface{}) {
	q.DelayQueue.AddAfter(item, q.rateLimiter.When())
}

func (q *rateLimitingQueue) Close() {
	q.DelayQueue.Close()
}

func NewRateLimitingQueue(opts RateLimitingQueueOption) RateLimiterQueue {
	return &rateLimitingQueue{
		DelayQueue:  NewDelayingQueue(NewQueue()),
		rateLimiter: newBucketRateLimiter(opts.rate, opts.buckets),
	}
}

// RateLimitingQueueOption 限速队列配置
type RateLimitingQueueOption struct {
	// rate 往桶里放Token的速率
	rate float64
	// buckets token 桶的容量大小
	buckets int
}

// rateLimitingQueue 限速队列实现对象
type rateLimitingQueue struct {
	// DelayQueue 继承 DelayQueue Interface 延迟队列基本功能
	DelayQueue
	// rateLimiter 限速器
	rateLimiter RateLimiter
}

// RateLimiter 限速器对象
type RateLimiter interface {
	// When 获取需要多长时间才能入队
	When() time.Duration
}

// bucketRateLimiter 令牌桶限速器
type bucketRateLimiter struct {
	*rate.Limiter
}

func newBucketRateLimiter(limit float64, tokens int) *bucketRateLimiter {
	return &bucketRateLimiter{rate.NewLimiter(rate.Limit(limit), tokens)}
}

var _ RateLimiter = &bucketRateLimiter{}

func (r *bucketRateLimiter) When() time.Duration {
	// 通过 r.Limiter.Reserve().Delay 函数返回指定元素应该等待的时间
	return r.Limiter.Reserve().Delay()
}
