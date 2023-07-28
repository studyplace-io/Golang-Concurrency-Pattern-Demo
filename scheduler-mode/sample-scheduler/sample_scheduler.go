package sample_scheduler

import (
	"fmt"
	"sync"
)

// task 任务
type task func()

// scheduler 调度器
type scheduler struct {
	tasks   chan task     // 调度队列
	workers int           // 并发调度数量
	stopC   chan struct{} // 退出通知
	wg      sync.WaitGroup
}

func newScheduler(workers int) *scheduler {
	return &scheduler{
		workers: workers,
		tasks:   make(chan task, 2),
		stopC:   make(chan struct{}),
	}
}

// start 启动scheduler调度器
func (s *scheduler) start() {
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			for {
				// 从chan task拿出任务调度
				select {
				case <-s.stopC:
					fmt.Println("scheduler finished")
					return
				case t := <-s.tasks:
					t()
				}
			}
		}()
	}
}

// addTask 任务入队
func (s *scheduler) addTask(t task) {
	s.tasks <- t
}

// stop 调度器退出
func (s *scheduler) stop() {
	go func() {
		s.wg.Wait()
		defer func() {
			s.stopC <- struct{}{}
		}()
	}()

}
