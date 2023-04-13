package task_job_mode

import (
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/task-job-mode/model"
	"log"
	"sync"
)

// 默认最大并行数
const defaultWorkerNumLimit = 8

// ClientWorker 用于分发执行Server发来的请求。内部维护任务队列和任务执行协程池。
// 同时负责向数据库同步任务信息
type ClientWorker struct {
	queue           *taskQueue                       // 任务队列
	workerChan      chan struct{}                    // 控制并发数用的channel
	responseWriter  func(result *model.Result) error // 向Server响应结果用
	workerWaitGroup *sync.WaitGroup
	closed          bool
}

var pool *ClientWorker
var lock sync.Mutex

// CreateClientWorker 创建clientWorker
func CreateClientWorker(responseWriter func(result *model.Result) error) *ClientWorker {
	lock.Lock()
	defer lock.Unlock()
	if pool == nil {
		log.Printf("Create new worker pool...")
		var workerLimit int
		//workerLimit := common.AppConfig.Worker.ConcurrencyLimit
		if workerLimit <= 0 {
			workerLimit = defaultWorkerNumLimit
		}
		pool = &ClientWorker{
			responseWriter:  responseWriter,
			workerChan:      make(chan struct{}, workerLimit),
			queue:           &taskQueue{},
			workerWaitGroup: &sync.WaitGroup{},
		}
	}
	return pool
}

// Execute 接收从server来的新任务。完成任务存储后返回，分别存储至MySQL和内存队列中。
// 此方法返回仅代表任务接收，并不代表任务执行完成。任务是异步执行的
func (c *ClientWorker) Execute(task *model.Task) (err error) {
	// 关闭中时拦截任务
	if c.closed {
		return AlreadyClosed
	}
	// 构建任务信息
	log.Printf("Worker pool receive new task: %s", task.TaskId)

	node := &taskNode{task: task}
	//// 新任务写入数据库
	//if node.id, node.jobIds, err = dal.CreateTask(nil, task); err != nil {
	// if err == dal.TaskAlreadyExists {
	//  return nil
	// }
	// log.Errorf("Save task into to database failed: %v", err)
	// return
	//}

	// 将任务写入内存的队列中，并调度执行
	select {
	case c.workerChan <- struct{}{}:
		// 并发量没到上限，直接调度执行
		go func() {
			c.workerWaitGroup.Add(1)
			defer c.workerWaitGroup.Done()
			executeTask(node)
		}()
	default:
		// 写入Queue
		c.queue.push(node)
	}
	return
}

// Close 关闭ClientWorker。执行此方法后将不再接收新任务，
// Close方法会在将现有任务执行完成后返回
func (c *ClientWorker) Close() error {
	lock.Lock()
	defer lock.Unlock()
	// 标识当前worker关闭
	c.closed = true

	// 无活跃worker pool则直接跳过
	if pool == nil {
		log.Printf("No active worker pool exists")
		return nil
	}
	// 关闭worker pool
	log.Printf("Worker pool closing...")
	close(c.workerChan)
	// 等待已有任务执行完成
	c.workerWaitGroup.Wait()
	pool = nil
	return nil
}
