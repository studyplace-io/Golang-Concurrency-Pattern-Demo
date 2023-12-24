package worker

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	worker := NewWorker()

	worker.RegisterJob("test1")
	worker.RegisterJob("test2-with-context")

	err := worker.RunCronJob("test1", time.Second, func() {
		fmt.Println("run job1...")
	})

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	err = worker.RunCronJobWithContext("test2-with-context", ctx, time.Second, func(c context.Context) {
		select {
		case <-c.Done():
			fmt.Println("timeout...")
		default:
		}
		fmt.Println("run job2...")
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(worker.JobStatus("test1"))

	// 停止
	select {
	case <-time.After(time.Second * 20):
		worker.StopJob("test1")
		worker.StopJob("test2-with-context")
	}

}
