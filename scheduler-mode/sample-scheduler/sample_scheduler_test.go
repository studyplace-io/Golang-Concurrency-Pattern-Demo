package sample_scheduler

import (
	"fmt"
	"testing"
	"time"
)

func TestSampleScheduler(test *testing.T) {
	s := newScheduler(2)
	s.start()
	defer func() {
		s.stop()
	}()

	s.addTask(func() {
		fmt.Println("Task 1")
	})

	s.addTask(func() {
		fmt.Println("Task 2")
	})
	time.Sleep(time.Second * 2)
	s.addTask(func() {
		fmt.Println("Task 3")
	})

	s.addTask(func() {
		fmt.Println("Task 3")
	})

	s.addTask(func() {
		fmt.Println("Task 3")
	})

}
